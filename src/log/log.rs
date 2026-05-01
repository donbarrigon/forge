use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Log {
    pub time: String,
    pub level: String,
    pub message: String,
    pub line: String,
    pub file: String,
    pub data: Option<Value>,
}

pub fn emergency(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn alert(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn critical(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn error(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn warning(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn notice(message: impl Into<String>, data: Option<Value>) {
    //
}

pub fn info(message: impl Into<String>, data: Option<Value>) {
    //
}
