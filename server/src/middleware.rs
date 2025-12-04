use axum::{
    extract::{ConnectInfo, Request, State},
    http::{header, StatusCode},
    middleware::Next,
    response::{IntoResponse, Response},
};
use jsonwebtoken::{decode, DecodingKey, Validation};
use std::net::SocketAddr;

use crate::config::get_jwt_secret;
use crate::state::AppState;
use crate::types::Claims;

pub async fn auth_middleware(
    headers: axum::http::HeaderMap,
    request: Request,
    next: Next,
) -> Response {
    if let Some(auth) = headers
        .get(header::AUTHORIZATION)
        .and_then(|h| h.to_str().ok())
    {
        if let Some(token) = auth.strip_prefix("Bearer ") {
            if decode::<Claims>(
                token,
                &DecodingKey::from_secret(get_jwt_secret().as_bytes()),
                &Validation::default(),
            )
            .is_ok()
            {
                return next.run(request).await;
            }
        }
    }
    StatusCode::UNAUTHORIZED.into_response()
}

/// Basic CC mitigation: per-IP sliding window rate limiter across HTTP and WebSocket handshakes.
pub async fn rate_limit_middleware(
    State(state): State<AppState>,
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    request: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    let client_ip = addr.ip().to_string();
    if !state.rate_limiter.is_allowed(&client_ip).await {
        tracing::warn!("Rate limit exceeded for {}", client_ip);
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    Ok(next.run(request).await)
}
