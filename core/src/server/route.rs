use crate::server::controller::Controller;
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct RouteInfo {
    pub method: String,
    pub path: String,
}

pub type RouteMap = HashMap<String, RouteInfo>;

#[derive(Clone)]
pub struct Route {
    pub path: &'static str,
    pub controller: Option<Controller>,
    pub children: Vec<Route>,
    pub map: RouteMap,
    pub is_dinamic: bool,
}

impl Route {
    pub fn new() -> Self {
        return Self {
            path: "",
            controller: None,
            children: Vec::new(),
            map: RouteMap::new(),
            is_dinamic: false,
        };
    }
}
