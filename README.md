# Minefield

## Author

Sam Choukri

```
email: sam@choukri.net
phone: 858-382-3062
```

## Challenge

Write a program which takes as input a list of mines composing a 2D minefield; each mine has an X position, a Y position, and an explosive power. All three parameters may be assumed to be single-precision floats; explosive power may not be negative. There may not be more than one mine at the same coordinates.
 
When a mine in the minefield is triggered at time T=0, it causes all other mines within a straight-line distance less than or equal to its explosive power to be triggered at time T=1. Those mines subsequently trigger additional mines at T=2, and so forth, in a chain reaction.
 
Have your program determine, for any given input minefield, the mine that, if triggered first, will result in the highest number of explosions occurring during a single time interval. Output the coordinates of the winning mine, the time interval of the peak number of explosions, and the number of explosions during that interval. In case of a tie, output each of the best mines, sorted by X coordinate then Y coordinate.
 
Assume that the minefield may be large, but not larger than can easily fit in memory; optimize for processing efficiency.

## Assumptions

It is possible for a given starting mine to have more than one time interval with the same peak number of explosions. It was not specified in the instructions which interval to use in case of a tie, but I chose to use the earliest time interval. This edge-case is covered in my unit tests.

The minefield application takes as input a plain-text file with 3 whitespace-delimited (tabs or spaces) fields per line:

1. X coordinate (positive or negative float32)
2. Y coordinate (positive or negative float32)
3. Explosive Power (positive float32)

The applicaiton will abort with an error message if it encounters any invalid data. I have provided several valid input files that I used for my own testing:

```
input.txt       # 1,000 mines (default)
input_small.txt # 10 mines
input_large.txt # 10,000 mines
```


## Running Minefield

minefield is written in Go. You can build and run it directly from the source files if you have a working Go environment installed. Or you can build a docker image and run minefield inside a container.

### Running with Go

I used Go version 1.11.4 on a Mac to develop and run the `minefield` application. It should work fine with earlier versions of Go, but I have not tested that assumption.

Move or copy the supplied `minefield` directory somewhere inside your `$GOPATH`:

```
mv minefield $GOPATH/src/
```

Change directories into the `minefield` directory:

```
cd $GOPATH/src/minefield
```

Run the tests to make sure everything passes:

```
go test -v
```

Build the `minefield` application:

```
go build -o minefield
```

That will produce a `minefield` application inside of the same directory:


Run the `minefield` application with the default input file `input.txt`:

```
./minefield
```

The output should be:

```
Winner (0): Mine ID=263, X=2.104142, Y=0.367990, Peak Time=4, Peak Explosions=306
Winner (1): Mine ID=326, X=3.200606, Y=-0.482302, Peak Time=5, Peak Explosions=306
```

If you want to specify a different input file, use the `-file` command line flag:

```
./minefield -file <FILEPATH>
```

### Running with Docker

To run the `minefield` application using Docker, you will need a recent version of Docker installed and ready to use. 

First, change directories to the supplied `minefield` directory:

```
cd minefield
```

Build the docker image and give it whatever name you want (I'm going to name it `sam-minefield`):

```
docker build -t sam-minefield .
```

This docker image uses the official Go 1.11 image. If you have not previously used this image on your computer, it will take a minute or two to download all the data.


Once the image is done building, you can run the `minefield` application inside an ephemeral docker container like this:

```
docker run --rm sam-minefield
```

The output should be the same as described earlier:

```
Winner (0): Mine ID=263, X=2.104142, Y=0.367990, Peak Time=4, Peak Explosions=306
Winner (1): Mine ID=326, X=3.200606, Y=-0.482302, Peak Time=5, Peak Explosions=306
```

To use a different input file, you need to specify a volume for the docker run command. The volume maps a directory on your host computer to a directory inside the container. The easiest thing to do is to just map your current working directory `$PWD` to the working directory inside of the container. Make sure the file you want to specify is inside your current working directory, then run docker like this:

```
docker run -v $PWD:/data -w /data --rm sam-minefield -file your_input_file.txt
```

## Benchmarks

Here are some timings for running `minefield` on my Mac with the 3 different size input files supplied:

`input.txt` (1,000 mines, the default)

```
time ./minefield -file input.txt 
Winner (0): Mine ID=263, X=2.104142, Y=0.367990, Peak Time=4, Peak Explosions=306
Winner (1): Mine ID=326, X=3.200606, Y=-0.482302, Peak Time=5, Peak Explosions=306

real    0m0.161s
user    0m0.173s
sys     0m0.011s
```

`input_small.txt` (10 mines)

```
time ./minefield -file input_small.txt 
Winner (0): Mine ID=1, X=0.253293, Y=-0.591300, Peak Time=1, Peak Explosions=9

real    0m0.009s
user    0m0.003s
sys     0m0.004s
```

`input_large.txt` (10,000 mines)


```
time ./minefield -file input_large.txt 
Winner (0): Mine ID=6342, X=-5.933625, Y=-0.270309, Peak Time=21, Peak Explosions=737
Winner (1): Mine ID=6823, X=-5.930994, Y=-0.143614, Peak Time=22, Peak Explosions=737
Winner (2): Mine ID=7664, X=-5.809915, Y=-0.216525, Peak Time=22, Peak Explosions=737

real    0m13.005s
user    0m25.772s
sys     0m0.531s
