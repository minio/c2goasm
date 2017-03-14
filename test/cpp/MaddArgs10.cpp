#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

void MaddArgs10(float* arg1, float* arg2, float* arg3, float* arg4, float* arg5, float* arg6, float* arg7, float* arg8, float* arg9, float* result) {
    __m256 vec1 = _mm256_load_ps(arg1);
    __m256 vec2 = _mm256_load_ps(arg2);
    __m256 vec3 = _mm256_load_ps(arg3);
    __m256 vec4 = _mm256_load_ps(arg4);
    __m256 vec5 = _mm256_load_ps(arg5);
    __m256 vec6 = _mm256_load_ps(arg6);
    __m256 vec7 = _mm256_load_ps(arg7);
    __m256 vec8 = _mm256_load_ps(arg8);
    __m256 vec9 = _mm256_load_ps(arg9);
    __m256 res1  = _mm256_fmadd_ps(vec1, vec2, vec3);
    __m256 res2  = _mm256_fmadd_ps(res1, vec4, vec5);
    __m256 res3  = _mm256_fmadd_ps(res2, vec6, vec7);
    __m256 res4  = _mm256_fmadd_ps(res3, vec8, vec9);
    _mm256_storeu_ps(result, res4);
}
