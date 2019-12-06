use std::env;
use std::fs::File;
use std::io::{self, prelude::*};

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() != 3 {
        println!("Must supply two arguments (input file location, desired number)");
        std::process::exit(1)
    }

    let input_path = &args[1];
    let desired_num: usize = (&args[2]).parse().unwrap();
    let mut file = File::open(input_path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;

    let intcode: Vec<usize> = contents
        .trim()
        .split(",")
        .map(|x| x.parse().unwrap())
        .collect();

    match test_intcode(intcode, desired_num) {
        Some((n, v)) => println!("n = {}, v = {}, 100 * n + v = {}", n, v, 100 * n + v),
        None => println!("Failed to find combination"),
    }
    Ok(())
}

fn test_intcode(initial_state: Vec<usize>, desired_num: usize) -> Option<(usize, usize)> {
    for n in 0..=99 {
        for v in 0..=99 {
            let mut intcode = initial_state.clone();
            intcode[1] = n;
            intcode[2] = v;
            intcode = run_intcode(intcode);
            if intcode[0] == desired_num {
                return Some((n, v));
            }
        }
    }
    None
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
