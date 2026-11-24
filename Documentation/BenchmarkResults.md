## Serial Benchmark
```text
goos: darwin
goarch: amd64
pkg: boids/tests
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkSerial/25-8         	     100	     18418 ns/op
BenchmarkSerial/25-8         	     100	     11832 ns/op
BenchmarkSerial/25-8         	     100	     12257 ns/op
BenchmarkSerial/25-8         	     100	     12856 ns/op
BenchmarkSerial/25-8         	     100	     11866 ns/op
BenchmarkSerial/50-8         	     100	     70414 ns/op
BenchmarkSerial/50-8         	     100	     58417 ns/op
BenchmarkSerial/50-8         	     100	     74483 ns/op
BenchmarkSerial/50-8         	     100	     74259 ns/op
BenchmarkSerial/50-8         	     100	     66015 ns/op
BenchmarkSerial/100-8        	     100	    219541 ns/op
BenchmarkSerial/100-8        	     100	    208959 ns/op
BenchmarkSerial/100-8        	     100	    207789 ns/op
BenchmarkSerial/100-8        	     100	    204542 ns/op
BenchmarkSerial/100-8        	     100	    206559 ns/op
BenchmarkSerial/500-8        	     100	   5068479 ns/op
BenchmarkSerial/500-8        	     100	   5209285 ns/op
BenchmarkSerial/500-8        	     100	   5158066 ns/op
BenchmarkSerial/500-8        	     100	   5051282 ns/op
BenchmarkSerial/500-8        	     100	   5204094 ns/op
BenchmarkSerial/1000-8       	     100	  20596590 ns/op
BenchmarkSerial/1000-8       	     100	  21225115 ns/op
BenchmarkSerial/1000-8       	     100	  21814926 ns/op
BenchmarkSerial/1000-8       	     100	  22366155 ns/op
BenchmarkSerial/1000-8       	     100	  21085500 ns/op
PASS
ok  	boids/tests	13.885s
```

## Parallel Benchmark
```text
goos: darwin
goarch: amd64
pkg: boids/tests
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkParallel/25/1_workers-8         	     100	     25944 ns/op
BenchmarkParallel/25/1_workers-8         	     100	     17755 ns/op
BenchmarkParallel/25/1_workers-8         	     100	     19578 ns/op
BenchmarkParallel/25/1_workers-8         	     100	     18223 ns/op
BenchmarkParallel/25/1_workers-8         	     100	     17710 ns/op
BenchmarkParallel/25/2_workers-8         	     100	     12308 ns/op
BenchmarkParallel/25/2_workers-8         	     100	     12598 ns/op
BenchmarkParallel/25/2_workers-8         	     100	     12232 ns/op
BenchmarkParallel/25/2_workers-8         	     100	     13800 ns/op
BenchmarkParallel/25/2_workers-8         	     100	     16199 ns/op
BenchmarkParallel/25/4_workers-8         	     100	     15265 ns/op
BenchmarkParallel/25/4_workers-8         	     100	     16077 ns/op
BenchmarkParallel/25/4_workers-8         	     100	     15186 ns/op
BenchmarkParallel/25/4_workers-8         	     100	     13862 ns/op
BenchmarkParallel/25/4_workers-8         	     100	     15062 ns/op
BenchmarkParallel/25/8_workers-8         	     100	     16486 ns/op
BenchmarkParallel/25/8_workers-8         	     100	     21868 ns/op
BenchmarkParallel/25/8_workers-8         	     100	     19260 ns/op
BenchmarkParallel/25/8_workers-8         	     100	     15063 ns/op
BenchmarkParallel/25/8_workers-8         	     100	     15117 ns/op
BenchmarkParallel/25/10_workers-8        	     100	     15718 ns/op
BenchmarkParallel/25/10_workers-8        	     100	     15810 ns/op
BenchmarkParallel/25/10_workers-8        	     100	     17016 ns/op
BenchmarkParallel/25/10_workers-8        	     100	     15614 ns/op
BenchmarkParallel/25/10_workers-8        	     100	     15623 ns/op
BenchmarkParallel/50/1_workers-8         	     100	     69794 ns/op
BenchmarkParallel/50/1_workers-8         	     100	     65167 ns/op
BenchmarkParallel/50/1_workers-8         	     100	     66396 ns/op
BenchmarkParallel/50/1_workers-8         	     100	     62212 ns/op
BenchmarkParallel/50/1_workers-8         	     100	     62733 ns/op
BenchmarkParallel/50/2_workers-8         	     100	     39318 ns/op
BenchmarkParallel/50/2_workers-8         	     100	     39132 ns/op
BenchmarkParallel/50/2_workers-8         	     100	     38362 ns/op
BenchmarkParallel/50/2_workers-8         	     100	     39050 ns/op
BenchmarkParallel/50/2_workers-8         	     100	     38474 ns/op
BenchmarkParallel/50/4_workers-8         	     100	     32372 ns/op
BenchmarkParallel/50/4_workers-8         	     100	     31629 ns/op
BenchmarkParallel/50/4_workers-8         	     100	     29836 ns/op
BenchmarkParallel/50/4_workers-8         	     100	     30035 ns/op
BenchmarkParallel/50/4_workers-8         	     100	     30897 ns/op
BenchmarkParallel/50/8_workers-8         	     100	     34405 ns/op
BenchmarkParallel/50/8_workers-8         	     100	     33300 ns/op
BenchmarkParallel/50/8_workers-8         	     100	     33957 ns/op
BenchmarkParallel/50/8_workers-8         	     100	     39462 ns/op
BenchmarkParallel/50/8_workers-8         	     100	     36568 ns/op
BenchmarkParallel/50/10_workers-8        	     100	     34742 ns/op
BenchmarkParallel/50/10_workers-8        	     100	     35371 ns/op
BenchmarkParallel/50/10_workers-8        	     100	     35613 ns/op
BenchmarkParallel/50/10_workers-8        	     100	     41228 ns/op
BenchmarkParallel/50/10_workers-8        	     100	     36279 ns/op
BenchmarkParallel/100/1_workers-8        	     100	    236107 ns/op
BenchmarkParallel/100/1_workers-8        	     100	    235158 ns/op
BenchmarkParallel/100/1_workers-8        	     100	    238148 ns/op
BenchmarkParallel/100/1_workers-8        	     100	    234269 ns/op
BenchmarkParallel/100/1_workers-8        	     100	    230522 ns/op
BenchmarkParallel/100/2_workers-8        	     100	    148640 ns/op
BenchmarkParallel/100/2_workers-8        	     100	    148274 ns/op
BenchmarkParallel/100/2_workers-8        	     100	    146447 ns/op
BenchmarkParallel/100/2_workers-8        	     100	    149318 ns/op
BenchmarkParallel/100/2_workers-8        	     100	    144002 ns/op
BenchmarkParallel/100/4_workers-8        	     100	    107552 ns/op
BenchmarkParallel/100/4_workers-8        	     100	    107595 ns/op
BenchmarkParallel/100/4_workers-8        	     100	    112351 ns/op
BenchmarkParallel/100/4_workers-8        	     100	    109791 ns/op
BenchmarkParallel/100/4_workers-8        	     100	    117916 ns/op
BenchmarkParallel/100/8_workers-8        	     100	    116086 ns/op
BenchmarkParallel/100/8_workers-8        	     100	    112463 ns/op
BenchmarkParallel/100/8_workers-8        	     100	    125931 ns/op
BenchmarkParallel/100/8_workers-8        	     100	    179347 ns/op
BenchmarkParallel/100/8_workers-8        	     100	    155403 ns/op
BenchmarkParallel/100/10_workers-8       	     100	    128143 ns/op
BenchmarkParallel/100/10_workers-8       	     100	    125074 ns/op
BenchmarkParallel/100/10_workers-8       	     100	    121197 ns/op
BenchmarkParallel/100/10_workers-8       	     100	    116821 ns/op
BenchmarkParallel/100/10_workers-8       	     100	    124520 ns/op
BenchmarkParallel/500/1_workers-8        	     100	   5824954 ns/op
BenchmarkParallel/500/1_workers-8        	     100	   5227752 ns/op
BenchmarkParallel/500/1_workers-8        	     100	   5568794 ns/op
BenchmarkParallel/500/1_workers-8        	     100	   6397593 ns/op
BenchmarkParallel/500/1_workers-8        	     100	   6394390 ns/op
BenchmarkParallel/500/2_workers-8        	     100	   3355630 ns/op
BenchmarkParallel/500/2_workers-8        	     100	   3321165 ns/op
BenchmarkParallel/500/2_workers-8        	     100	   3503515 ns/op
BenchmarkParallel/500/2_workers-8        	     100	   3635007 ns/op
BenchmarkParallel/500/2_workers-8        	     100	   3377480 ns/op
BenchmarkParallel/500/4_workers-8        	     100	   2884372 ns/op
BenchmarkParallel/500/4_workers-8        	     100	   3049441 ns/op
BenchmarkParallel/500/4_workers-8        	     100	   2915469 ns/op
BenchmarkParallel/500/4_workers-8        	     100	   3276647 ns/op
BenchmarkParallel/500/4_workers-8        	     100	   2947194 ns/op
BenchmarkParallel/500/8_workers-8        	     100	   3339956 ns/op
BenchmarkParallel/500/8_workers-8        	     100	   3819216 ns/op
BenchmarkParallel/500/8_workers-8        	     100	   3407847 ns/op
BenchmarkParallel/500/8_workers-8        	     100	   3719136 ns/op
BenchmarkParallel/500/8_workers-8        	     100	   2995335 ns/op
BenchmarkParallel/500/10_workers-8       	     100	   2876566 ns/op
BenchmarkParallel/500/10_workers-8       	     100	   3160518 ns/op
BenchmarkParallel/500/10_workers-8       	     100	   3473183 ns/op
BenchmarkParallel/500/10_workers-8       	     100	   3663191 ns/op
BenchmarkParallel/500/10_workers-8       	     100	   3186726 ns/op
BenchmarkParallel/1000/1_workers-8       	     100	  24238271 ns/op
BenchmarkParallel/1000/1_workers-8       	     100	  24078982 ns/op
BenchmarkParallel/1000/1_workers-8       	     100	  21482829 ns/op
BenchmarkParallel/1000/1_workers-8       	     100	  21353250 ns/op
BenchmarkParallel/1000/1_workers-8       	     100	  24483561 ns/op
BenchmarkParallel/1000/2_workers-8       	     100	  13716092 ns/op
BenchmarkParallel/1000/2_workers-8       	     100	  15285183 ns/op
BenchmarkParallel/1000/2_workers-8       	     100	  15039266 ns/op
BenchmarkParallel/1000/2_workers-8       	     100	  14066224 ns/op
BenchmarkParallel/1000/2_workers-8       	     100	  12946935 ns/op
BenchmarkParallel/1000/4_workers-8       	     100	  11085078 ns/op
BenchmarkParallel/1000/4_workers-8       	     100	  10715641 ns/op
BenchmarkParallel/1000/4_workers-8       	     100	  11520591 ns/op
BenchmarkParallel/1000/4_workers-8       	     100	  11190008 ns/op
BenchmarkParallel/1000/4_workers-8       	     100	  11921778 ns/op
BenchmarkParallel/1000/8_workers-8       	     100	  15975843 ns/op
BenchmarkParallel/1000/8_workers-8       	     100	  16092274 ns/op
BenchmarkParallel/1000/8_workers-8       	     100	  15838066 ns/op
BenchmarkParallel/1000/8_workers-8       	     100	  13786920 ns/op
BenchmarkParallel/1000/8_workers-8       	     100	  13843754 ns/op
BenchmarkParallel/1000/10_workers-8      	     100	  13480924 ns/op
BenchmarkParallel/1000/10_workers-8      	     100	  12636252 ns/op
BenchmarkParallel/1000/10_workers-8      	     100	  13527348 ns/op
BenchmarkParallel/1000/10_workers-8      	     100	  15917664 ns/op
BenchmarkParallel/1000/10_workers-8      	     100	  15557735 ns/op
PASS
ok  	boids/tests	50.412s
```