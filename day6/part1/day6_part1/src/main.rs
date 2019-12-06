use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::io::{self, prelude::*, BufReader};

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        println!("Must supply one argument (input file location)");
        std::process::exit(1)
    }

    let input_path = &args[1];
    let file = File::open(input_path)?;
    let reader = BufReader::new(file);

    let mut orbits: HashMap<String, String> = HashMap::new();
    reader
        .lines()
        .map(|l| l.unwrap().split(")").map(|s| s.to_string()).collect())
        .for_each(|f: Vec<String>| {
            let _ = orbits.insert(f[1].to_owned(), f[0].to_owned());
        });

    let mut distances: HashMap<String, usize> = HashMap::new();
    for child in orbits.keys() {
        let mut distance = 0;
        let mut temp_child = child.clone();
        while temp_child != "COM" {
            distance += 1;
            temp_child = orbits.get(&temp_child).unwrap().to_string();
        }
        let _ = distances.insert(child.to_string(), distance);
    }

    let mut answer = 0;
    for distance in distances.values() {
        answer += distance;
    }

    println!("{:?}", answer);
    Ok(())
}
