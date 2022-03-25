use aes_gcm::{Aes128Gcm, Key, Nonce};
use aes_gcm::aead::{Aead, NewAead};
use getrandom::getrandom;

use crate::database::{EncryptedPasteModel, PasteModel};

pub fn generate_key() -> [u8; 16] {
    let mut key_bytes = [0; 16];
    getrandom(&mut key_bytes).unwrap();
    key_bytes
}

pub fn encrypt_paste_model(key_bytes: &[u8], model: PasteModel) -> Result<EncryptedPasteModel, aes_gcm::Error> {
    let key = Key::from_slice(key_bytes);
    let cipher = Aes128Gcm::new(key);

    let mut nonce_bytes = [0; 12];
    getrandom(&mut nonce_bytes).unwrap();
    let nonce = Nonce::from_slice(&nonce_bytes);

    let cipher_content = cipher.encrypt(nonce, model.content.as_bytes())?;

    Ok(EncryptedPasteModel {
        language: model.language,
        created_at: model.created_at,
        nonce: nonce_bytes,
        content: cipher_content,
    })
}

pub fn decrypt_paste_model(key_bytes: &[u8], model: EncryptedPasteModel) -> Result<PasteModel, aes_gcm::Error> {
    let key = Key::from_slice(key_bytes);
    let cipher = Aes128Gcm::new(key);

    let nonce = Nonce::from_slice(&model.nonce);

    let plain_content = cipher.decrypt(nonce, model.content.as_slice())?;
    let plain_text = String::from_utf8(plain_content).unwrap();

    Ok(PasteModel {
        language: model.language,
        created_at: model.created_at,
        content: plain_text,
    })
}