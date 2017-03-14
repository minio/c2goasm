#include <immintrin.h>

#define ALIGN(x) x __attribute__((aligned(32)))

static const ALIGN(float a[8]) = {1.0f, 2.0f, 3.0f, 4.0f, 5.0f, 6.0f, 7.0f, 8.0f};

void MaddConstant(float* arg1, float* arg2, float* result) {
    __m256 vec1 = _mm256_load_ps(arg1);
    __m256 vec2 = _mm256_load_ps(arg2);
    __m256 vec3 = _mm256_load_ps(a);
    __m256 res  = _mm256_fmadd_ps(vec1, vec2, vec3);
    _mm256_storeu_ps(result, res);
}
