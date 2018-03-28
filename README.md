# adcalc
output format: 
[filename] [doselimit] [ad] [ad30] [peakdose] [dose integral]

options:
   -D[value] : Limit dose for normal tissue.
   -N[value] : Density[ppm] for Normal tissue.
   -T[value] : Density[ppm] for Tumor tissue.
   -R[value] : T/N Ratio.
   
ex.:
  ./adcalc sample.zdis -D11.0 -N10 -T70
   means, Limit dose for normal tissue is 11Gy, boron for Normal tissue is 10ppm, and boron for Tumor tissue is 100ppm
