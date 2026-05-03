use crate::server::RouteMap;
use hyper::Request;
use hyper::body::Incoming;
use std::sync::Arc;

pub struct Context {
    pub req: Request<Incoming>,
    pub params: Vec<(String, String)>,
    pub route_map: Arc<RouteMap>,
}
