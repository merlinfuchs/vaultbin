use std::borrow::Cow;
use std::path::Path;

use actix_web::{get, HttpResponse, Responder, web};
use actix_web::web::Bytes;
use mime::Mime;
use rust_embed::RustEmbed;

#[derive(RustEmbed)]
#[folder = "$CARGO_MANIFEST_DIR/frontend/build"]
struct FrontendFiles;

fn get_mime_type_for_file(path: &Path) -> Mime {
    match path.extension().and_then(|v| v.to_str()) {
        Some(v) => match v {
            "html" => mime::TEXT_HTML,
            "js" => mime::APPLICATION_JAVASCRIPT,
            "png" => mime::IMAGE_PNG,
            "css" => mime::TEXT_CSS,
            "svg" => mime::IMAGE_SVG,
            _ => mime::APPLICATION_OCTET_STREAM
        }
        None => mime::APPLICATION_OCTET_STREAM
    }
}

#[get("/{path:.*}")]
pub async fn route_serve_frontend(path: web::Path<String>) -> impl Responder {
    let mut path = path.into_inner();
    let mut file = FrontendFiles::get(&path);
    if file.is_none() {
        path = String::from("index.html");
        file = FrontendFiles::get(&path);
    }

    if let Some(file) = file {
        let mime_type = get_mime_type_for_file(Path::new(&path));
        let body: Bytes = match file.data {
            Cow::Borrowed(bytes) => bytes.into(),
            Cow::Owned(bytes) => bytes.into(),
        };
        HttpResponse::Ok()
            .append_header(("Content-Type", mime_type))
            .body(body)
    } else {
        HttpResponse::NotFound().finish()
    }
}
