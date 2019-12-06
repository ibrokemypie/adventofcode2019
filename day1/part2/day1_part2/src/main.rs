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

    let mut fuel: usize = 0;
    reader
        .lines()
        .map(|l| l.unwrap().parse().unwrap())
        .for_each(|f: usize| fuel += calc_fuel(f));

    println!("{:?}", fuel);

    Ok(())
}

fn calc_fuel(weight: usize) -> usize {
    let mut total_fuel = 0;

    let mut current_weight = weight / 3;

    while current_weight.checked_sub(2).is_some() {
        current_weight -= 2;
        total_fuel += current_weight;
        current_weight = current_weight / 3;
    }

    return total_fuel;
}
