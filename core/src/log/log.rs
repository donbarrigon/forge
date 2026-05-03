use chrono::Utc;
use serde_json::Value;
use std::fs::{self, OpenOptions};
use std::io::{BufWriter, Write};
use std::path::Path;
use std::sync::OnceLock;
use std::sync::atomic::{AtomicU8, Ordering};
use tokio::sync::mpsc;
use tokio::time::{Duration, sleep};

use crate::config::env;

// ─────────────────────────────────────────────
// LOG LEVEL GLOBAL (convertido desde Env)
// ─────────────────────────────────────────────

static LOG_LEVEL: AtomicU8 = AtomicU8::new(7);
static LOG_SENDER: OnceLock<mpsc::UnboundedSender<String>> = OnceLock::new();

pub fn get_log_level() -> u8 {
    return LOG_LEVEL.load(Ordering::Relaxed);
}

pub fn set_log_level(lv: &str) {
    let level = match lv {
        "emergency" => 0u8,
        "alert" => 1,
        "critical" => 2,
        "error" => 3,
        "warning" => 4,
        "notice" => 5,
        "info" => 6,
        _ => 7, // fallback debug
    };

    LOG_LEVEL.store(level, Ordering::Relaxed);
}

// ─────────────────────────────────────────────
// init
// ─────────────────────────────────────────────

pub fn init() {
    set_log_level(env().log.level.as_str());
    init_log_writer();
    start_cleanup_task();
}

pub fn init_log_writer() {
    let (tx, mut rx) = mpsc::unbounded_channel::<String>();
    LOG_SENDER.set(tx).ok();

    tokio::spawn(async move {
        let logs_dir = Path::new("./tmp/logs");
        let _ = fs::create_dir_all(logs_dir);

        while let Some(line) = rx.recv().await {
            let filename = format!("{}.log", chrono::Utc::now().format("%Y-%m-%d"));
            let filepath = logs_dir.join(filename);

            match OpenOptions::new().create(true).append(true).open(&filepath) {
                Ok(file) => {
                    let mut buf = BufWriter::new(file);
                    if writeln!(buf, "{}", line).is_err() {
                        eprintln!("[forge::log] failed to write log — disk full or permission denied");
                    }
                }
                Err(e) => {
                    eprintln!("[forge::log] failed to open log file: {}", e);
                }
            }
        }

        // el canal se cerró — todos los senders fueron dropeados
        eprintln!("[forge::log] writer task ended — no more logs will be written");
    });
}

// ─────────────────────────────────────────────
// TIMESTAMP
// ─────────────────────────────────────────────

fn get_now_iso() -> String {
    return Utc::now().to_rfc3339_opts(chrono::SecondsFormat::Millis, true);
}

// ─────────────────────────────────────────────
// CLEANUP TASK
// ─────────────────────────────────────────────

pub fn start_cleanup_task() {
    tokio::spawn(async {
        loop {
            cleanup_old_logs();
            sleep(Duration::from_secs(24 * 60 * 60)).await;
        }
    });
}

pub fn cleanup_old_logs() {
    let log_env = &env().log;
    let days = log_env.days as i64;
    let logs_dir = Path::new("./tmp/logs");

    if !logs_dir.exists() {
        return;
    }

    let cutoff = chrono::Utc::now().naive_utc().date() - chrono::Duration::days(days);

    let entries = match fs::read_dir(logs_dir) {
        Ok(e) => e,
        Err(_) => return,
    };

    for entry in entries.flatten() {
        let path = entry.path();
        if let Some(stem) = path.file_stem().and_then(|s| s.to_str()) {
            if let Ok(file_date) = chrono::NaiveDate::parse_from_str(stem, "%Y-%m-%d") {
                if file_date < cutoff {
                    let _ = fs::remove_file(&path);
                }
            }
        }
    }
}

// ─────────────────────────────────────────────
// CREATE LOG FUNCTIONS
// ─────────────────────────────────────────────

const RESET: &str = "\x1b[0m";

fn create_log(level: &'static str, message: String, data: Option<Value>) {
    let time = get_now_iso();

    let data_str = match &data {
        Some(d) => d.to_string(),
        None => "{}".to_string(),
    };

    let mut json = String::with_capacity(time.len() + level.len() + message.len() + data_str.len() + 64);

    json.push_str(r#"{"time":""#);
    json.push_str(&time);
    json.push_str(r#"","level":""#);
    json.push_str(level);
    json.push_str(r#"","message":""#);
    json.push_str(&message);
    json.push_str(r#"","data":"#);
    json.push_str(&data_str);
    json.push('}');

    write_to_file(json);

    // ─── consola — formato legible ────────────────────────────────────────
    if env().log.print {
        let color = match level {
            "emergency" => "\x1b[1;35m", // magenta brillante
            "alert" => "\x1b[1;31m",     // rojo brillante
            "critical" => "\x1b[1;33m",  // amarillo brillante
            "error" => "\x1b[31m",       // rojo
            "warning" => "\x1b[33m",     // amarillo
            "notice" => "\x1b[36m",      // cyan
            "info" => "\x1b[32m",        // verde
            "debug" => "\x1b[37m",       // blanco
            _ => "\x1b[0m",              // reset
        };
        let mut line = String::with_capacity(time.len() + level.len() + message.len() + 16);

        line.push_str(&time);
        line.push_str(" ");
        line.push_str(color); // color on
        line.push_str("[");
        line.push_str(level);
        line.push_str("]");
        line.push_str(RESET); // color off
        line.push_str(" ");
        line.push_str(&message);

        eprintln!("{}", line);

        if let Some(d) = &data {
            if let Ok(pretty) = serde_json::to_string_pretty(d) {
                eprintln!("{}{}{}", color, pretty, RESET);
            }
        }
    }
}

fn write_to_file(line: String) {
    match LOG_SENDER.get() {
        Some(tx) => {
            if tx.send(line).is_err() {
                eprintln!("[forge::log] writer channel closed — log lost");
            }
        }
        None => {
            eprintln!("[forge::log] writer not initialized — call init_log() first");
        }
    }
}

// ─────────────────────────────────────────────
// PUBLIC LOG FUNCTIONS
// ─────────────────────────────────────────────

pub fn emergency(msg: impl Into<String>, data: Option<Value>) {
    // emergency siempre se loggea — sin check de nivel
    create_log("emergency", msg.into(), data);
}

pub fn alert(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 1u8 {
        return;
    }
    create_log("alert", msg.into(), data);
}

pub fn critical(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 2u8 {
        return;
    }
    create_log("critical", msg.into(), data);
}

pub fn error(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 3u8 {
        return;
    }
    create_log("error", msg.into(), data);
}

pub fn warning(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 4u8 {
        return;
    }
    create_log("warning", msg.into(), data);
}

pub fn notice(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 5u8 {
        return;
    }
    create_log("notice", msg.into(), data);
}

pub fn info(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 6u8 {
        return;
    }
    create_log("info", msg.into(), data);
}

pub fn debug(msg: impl Into<String>, data: Option<Value>) {
    if get_log_level() < 7u8 {
        return;
    }
    create_log("debug", msg.into(), data);
}
