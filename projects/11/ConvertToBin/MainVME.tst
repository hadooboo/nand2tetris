// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/11/ProgramFlow/FibonacciSeries/FibonacciSeriesVME.tst

load Main.vm,
output-file Main.out,
output-list RAM[8001]%D1.6.2 RAM[8002]%D1.6.2 RAM[8003]%D1.6.2 RAM[8004]%D1.6.2
            RAM[8005]%D1.6.2 RAM[8006]%D1.6.2 RAM[8007]%D1.6.2 RAM[8008]%D1.6.2
            RAM[8009]%D1.6.2 RAM[8010]%D1.6.2 RAM[8011]%D1.6.2 RAM[8012]%D1.6.2
            RAM[8013]%D1.6.2 RAM[8014]%D1.6.2 RAM[8015]%D1.6.2 RAM[8016]%D1.6.2;

set RAM[8000] 1234,

repeat 100000 {
  vmstep;
}

output;
