use crate::errors::Error;
use crate::server::context::Context;
use http_body_util::Full;
use hyper::Response;
use hyper::body::Bytes;
use std::future::Future;
use std::pin::Pin;
use std::sync::Arc;

pub type ControllerFuture = Pin<Box<dyn Future<Output = Result<Response<Full<Bytes>>, Error>> + Send>>;
pub type Controller = Arc<dyn Fn(Context) -> ControllerFuture + Send + Sync>;
