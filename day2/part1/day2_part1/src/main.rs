use std::env;
use std::fs::File;
use std::io::{self, prelude::*};

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        println!("Must supply one argument (input file location)");
        std::process::exit(1)
    }

    let input_path = &args[1];
    let mut file = File::open(input_path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;

    let mut intcode: Vec<usize> = contents
        .trim()
        .split(",")
        .map(|x| x.parse().unwrap())
        .collect();

    intcode[1] = 12;
    intcode[2] = 2;

    intcode = run_intcode(intcode);
    println!("Value at position 0 is: {}", intcode[0]);
    Ok(())
}

fn run_intcode(mut intcode: Vec<usize>) -> Vec<usize> {
    let mut p = 0;

    while p < intcode.len() - 1 {
        let opcode = intcode[p];
        let val_one = intcode[intcode[p + 1]];
        let val_two = intcode[intcode[p + 2]];
        let output_pos = intcode[p + 3];

        match opcode {
            1 => {
                intcode[output_pos] = val_one + val_two;
            }
            2 => {
                intcode[output_pos] = val_one * val_two;
            }
            99 => break,
            _ => panic!(),
        }
        p += 4
    }

    return intcode;
}
