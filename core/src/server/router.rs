use crate::log;
use crate::server::RouteMap;
use crate::server::route::Route;
use http_body_util::Full;
use hyper::body::Bytes;
use hyper::{Request, Response};
use std::convert::Infallible;

#[derive(Clone)]
pub struct Router {
    pub route: Route,
    pub map: RouteMap,
}

impl Router {
    pub fn new() -> Self {
        return Self {
            route: Route::new(),
            map: RouteMap::new(),
        };
    }

    pub async fn handle(&self, req: Request<hyper::body::Incoming>) -> Result<Response<Full<Bytes>>, Infallible> {
        log::debug(format!("{}:{}", req.method(), req.uri().path()), None);
        return Ok(Response::new(Full::new(Bytes::from("Hello, from Forge!"))));
    }
}
