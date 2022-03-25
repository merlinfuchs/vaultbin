const FRONTEND_DIR: &str = concat!(env!("CARGO_MANIFEST_DIR"), "/frontend/build");

fn main() {
    // this forces the program to recompile when the frontend has changed
    println!("cargo:rerun-if-changed={}/src", FRONTEND_DIR)
}