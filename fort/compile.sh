
gfortran -c -O3 -fpic ./mo_par_DLS.f90
gfortran -c -O3 -fpic ./mo_DLS.f90
gfortran -c -O3 -fpic ./mo_alloc1.f90
gfortran -c -O3 -fpic ./mo_alloc.f90

gfortran -c -O3 -fpic *.f90

gfortran -shared -O3 -o libspheroid.so *.o

