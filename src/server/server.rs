use std::convert::Infallible;
use std::net::SocketAddr;

use http_body_util::Full;
use hyper::body::Bytes;
use hyper::service::service_fn;
use hyper::{Request, Response};
use hyper_util::rt::{TokioExecutor, TokioIo};
use hyper_util::server::conn::auto;
use tokio::net::TcpListener;
use tokio::sync::oneshot;

pub struct Server {
    addr: SocketAddr,
    shutdown_sender: Option<oneshot::Sender<()>>,
}

impl Server {
    pub fn new(host: &str, port: u16) -> Self {
        let addr = format!("{}:{}", host, port)
            .parse()
            .expect("dirección inválida");

        Self {
            addr,
            shutdown_sender: None,
        }
    }

    pub async fn listen(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        let listener = TcpListener::bind(self.addr).await?;

        // canal oneshot — manda una sola señal para apagar
        let (tx, mut rx) = oneshot::channel::<()>();
        self.shutdown_sender = Some(tx);

        println!("forge corriendo en http://{}", self.addr);

        loop {
            tokio::select! {
                result = listener.accept() => {
                    let (stream, _) = result?;
                    let io = TokioIo::new(stream);

                    tokio::spawn(async move {
                        if let Err(e) = auto::Builder::new(TokioExecutor::new())
                            .serve_connection(io, service_fn(hello))
                            .await
                        {
                            eprintln!("error en conexión: {:?}", e);
                        }
                    });
                }

                // señal manual de shutdown
                _ = &mut rx => {
                    println!("servidor {} detenido", self.addr);
                    break;
                }

                // Ctrl+C global
                _ = tokio::signal::ctrl_c() => {
                    println!("servidor {} detenido por Ctrl+C", self.addr);
                    break;
                }
            }
        }

        Ok(())
    }

    // apaga este servidor específico
    pub fn shutdown(&mut self) {
        if let Some(tx) = self.shutdown_sender.take() {
            let _ = tx.send(());
        }
    }
}

async fn hello(_req: Request<hyper::body::Incoming>) -> Result<Response<Full<Bytes>>, Infallible> {
    Ok(Response::new(Full::new(Bytes::from("Hello, from Forge!"))))
}
