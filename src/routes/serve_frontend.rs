use std::borrow::Cow;
use std::path::Path;

use actix_web::{get, HttpResponse, Responder, web};
use actix_web::web::Bytes;
use mime::Mime;
use rust_embed::RustEmbed;

use crate::database::Database;

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

fn cow_to_bytes(cow: Cow<'static, [u8]>) -> Bytes {
    match cow {
        Cow::Borrowed(bytes) => bytes.into(),
        Cow::Owned(bytes) => bytes.into(),
    }
}

// panics when there is no index.html
fn get_index_file(_path: &Path, _db: &Database) -> Bytes {
    return cow_to_bytes(FrontendFiles::get("index.html").unwrap().data);
    /* let maybe_paste_id = match path.file_name() {
        Some(n) if n.len() > 16 && n.len() < 32 => n.to_str(),
        _ => None
    };

    let fallback = || cow_to_bytes(FrontendFiles::get("index.html").unwrap().data);

    let paste_id = match maybe_paste_id {
        Some(p) => p,
        None => return fallback()
    };

    let decode_paste_id = match decode_bytes_from_string(&paste_id) {
        Ok(p) => p,
        Err(_) => return fallback()
    };

    let paste = match db.get_paste(&decode_paste_id) {
        Ok(Some(p)) => p,
        _ => return fallback()
    };

    let file = FrontendFiles::get("index.html").unwrap();
    let mut output = Vec::with_capacity(file.data.len());

    let mut rewriter = HtmlRewriter::new(
        lol_html::Settings {
            element_content_handlers: vec![
                element!("head", |el| {
                    let escaped_content = html_escape::encode_text(paste.content.unicode_truncate(500).0);
                    let description_meta = format!("<meta property=\"og:description\" content=\"{}\"/>", escaped_content);
                    el.append(&description_meta, ContentType::Html);
                    /* if let Some(language) = &paste.language {
                        let image_meta = format!("<meta property=\"og:image\" property=\"https://cdn.jsdelivr.net/npm/programming-languages-logos@0.0.3/src/{}/{}.png\"/>", language, language);
                        el.append(&image_meta, ContentType::Html);
                    } */
                    Ok(())
                })
            ],
            ..lol_html::Settings::default()
        },
        |c: &[u8]| output.extend_from_slice(c),
    );

    match rewriter.write(file.data.as_ref()) {
        Ok(_) => output.into(),
        // we really don't want to fail here
        Err(_) => cow_to_bytes(file.data)
    } */
}

#[get("/{path:.*}")]
pub async fn route_serve_frontend(path: web::Path<String>, db: web::Data<Database>) -> impl Responder {
    let raw_path = path.into_inner();
    let path = Path::new(&raw_path);

    let (body, mime_type) = match raw_path.as_str() {
        "index.html" => (get_index_file(&path, &db), mime::TEXT_HTML),
        p => match FrontendFiles::get(p) {
            Some(f) => (cow_to_bytes(f.data), get_mime_type_for_file(&path)),
            None => (get_index_file(&path, &db), mime::TEXT_HTML)
        }
    };

    HttpResponse::Ok()
        .append_header(("Content-Type", mime_type))
        .body(body)
}
