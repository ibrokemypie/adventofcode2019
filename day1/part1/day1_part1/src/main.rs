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

    let mut fuel: i32 = 0;
    reader
        .lines()
        .map(|l| l.unwrap().parse().unwrap())
        .for_each(|f: i32| fuel += f / 3 - 2);

    println!("{:?}", fuel);

    Ok(())
}
