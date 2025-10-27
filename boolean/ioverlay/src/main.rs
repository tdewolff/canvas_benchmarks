#![allow(warnings)]

use std::error::Error;
use std::fs::File;
use std::io::ErrorKind;
use std::io::BufReader;
use std::io::Write;
use std::fs::OpenOptions;
use std::time::Instant;

use serde::Serialize;

use i_overlay::core::fill_rule::FillRule;
use i_overlay::core::overlay_rule::OverlayRule;
use i_overlay::float::single::SingleFloatOverlay;

type Shape = Vec<Vec<[f64; 2]>>;

#[derive(Serialize)]
struct Res {
    Name: String,
    Z: usize,
    T: u128,
    Status: u64,
    Result: Vec<Shape>,
}

fn load_paths(name: &str) -> Result<Vec<Shape>, Box<dyn Error>>{
    let mut i = 0;
    let mut result = Vec::<Shape>::new();
    loop {
        let path = format!("../data/{}_{}.json", name, i);
        let file = match File::open(path) {
            Ok(x) => x,
            Err(e) => {
                if e.kind() == ErrorKind::NotFound {
                    return Ok(result);
                }
                return Err(Box::new(e));
            }
        };
        let reader = BufReader::new(file);
        result.push(serde_json::from_reader(reader)?);
        i += 1;
    }
}

fn main() {
    let europe = load_paths("europe").unwrap();
    let chile = load_paths("chile").unwrap();

    let mut file = OpenOptions::new()
        .write(true)
        .append(true)
        .open("../data/results.json")
        .unwrap();
    
    for z in 0..europe.len() {
        let now = Instant::now();
        let result = europe[z].overlay(&chile[z], OverlayRule::Union, FillRule::EvenOdd);
        for _ in 1..5 {
            _ = europe[z].overlay(&chile[z], OverlayRule::Union, FillRule::EvenOdd);
        }
        let elapsed = now.elapsed()/5;

        let json = Res{
            Name: "ioverlay".to_string(),
            Z: z,
            T: elapsed.as_nanos(),
            Status: 0,
            Result: result,
        };
        _ = serde_json::to_writer(&file, &json);
        _ = file.write(b"\n");
        println!("{:?} {:?}", z, elapsed);
    }
}
