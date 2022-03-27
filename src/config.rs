use std::path::Path;

use config::{Config, ConfigError, Environment, File};
use serde::{Deserialize, Serialize};

fn default_database_path() -> String {
    String::from("./data")
}

fn default_database_cache_size() -> u64 {
    100 * 1000 * 1000 // 100 MB
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DatabaseConfig {
    #[serde(default = "default_database_path")]
    pub path: String,
    #[serde(default = "default_database_cache_size")]
    pub cache_size: u64,
}

impl Default for DatabaseConfig {
    fn default() -> Self {
        Self {
            path: default_database_path(),
            cache_size: default_database_cache_size(),
        }
    }
}

fn default_ratelimit_burst_size() -> u32 {
    5
}

fn default_ratelimit_per_second() -> u64 {
    5
}

fn default_ratelimit_reverse_proxy() -> bool {
    false
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RatelimitConfig {
    #[serde(default = "default_ratelimit_burst_size")]
    pub burst_size: u32,
    #[serde(default = "default_ratelimit_per_second")]
    pub per_second: u64,
    #[serde(default = "default_ratelimit_reverse_proxy")]
    pub reverse_proxy: bool,
}

impl Default for RatelimitConfig {
    fn default() -> Self {
        Self {
            burst_size: default_ratelimit_burst_size(),
            per_second: default_ratelimit_per_second(),
            reverse_proxy: default_ratelimit_reverse_proxy()
        }
    }
}

fn default_host() -> String {
    "127.0.0.1".to_string()
}

fn default_port() -> u16 {
    8080
}

fn default_max_paste_size() -> usize {
    69420
}

fn default_max_expiration() -> u64 {
    60 * 60 * 24 * 30
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VaultbinConfig {
    #[serde(default)]
    pub database: DatabaseConfig,
    #[serde(default)]
    pub ratelimit: RatelimitConfig,

    #[serde(default = "default_host")]
    pub host: String,
    #[serde(default = "default_port")]
    pub port: u16,

    #[serde(default = "default_max_paste_size")]
    pub max_paste_size: usize,
    #[serde(default = "default_max_expiration")]
    pub max_expiration: u64
}

impl VaultbinConfig {
    pub fn new() -> Result<Self, ConfigError> {
        let mut config = Config::new();

        let config_file = "./Config.toml";
        if Path::new(config_file).exists() {
            config.merge(File::with_name(config_file))?;
        }

        config.merge(Environment::with_prefix("VAULTBIN").separator("__"))?;

        config.try_into()
    }
}
