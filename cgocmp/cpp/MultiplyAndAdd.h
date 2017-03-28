#include <immintrin.h>

void MultiplyAndAdd(float* arg1, float* arg2, float* arg3, float* result) {
    __m256 vec1 = _mm256_load_ps(arg1);
    __m256 vec2 = _mm256_load_ps(arg2);
    __m256 vec3 = _mm256_load_ps(arg3);
    __m256 res  = _mm256_fmadd_ps(vec1, vec2, vec3);
    _mm256_storeu_ps(result, res);
}
