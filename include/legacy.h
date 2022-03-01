#ifndef __LEGACY_H__
#define __LEGACY_H__


#define KRpar 25
#define KNpar 41
#define KMpar 181

#define  KN1par  41
#define  KR1par  25
#define  KM1par  181
#define  KREpar  35
#define  KIMpar  16
#define  KMD 2
#define  KSD 1 

extern int key, key_rd, keyel, keysub, keyls, key_org, key_fx, key_grid1, key_rd1, kn, km, kr, nratn, ndp;

extern float wl, rn, rk, pomin, pomax, xext, xabs, xsca, albedo, xblr, xldr;

extern float r[KRpar];
extern float grid[KNpar], sd[KNpar];
extern float rd[KRpar];


extern float f11[KMpar], f12[KMpar], f22[KMpar], f33[KMpar], f34[KMpar], f44[KMpar];
extern float angle[KMpar];
extern int key_sd, id, nmd, nsd;

extern float cm[KMD], sm[KMD], rmm[KMD];
extern float rrr[KNpar], ar[KNpar], xgrid[KNpar];
extern float ac;


void alloc_dls_array_(int *key, int *keyEL, int *keyAlloc);
	//void dls_read_input_();
void dls_read_config_(char *fname, int *len);
void optchar_(int *ndp);

void set_distname_o(char *fname);
void set_distname_f(char *fname);
void set_distname_n(char *fname);
void set_comm_name(int *i, char *s);
void sizedisdn_(int *, int *, int *, int *, float* , float* , float* , float* , float* , float* , float* , float *, int *, int* );
#endif

