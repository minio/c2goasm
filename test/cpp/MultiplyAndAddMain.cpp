#include <stdio.h>

void MultiplyAndAdd(float* arg1, float* arg2, float* arg3, float* result);

int main() {
    int i;
    float f1[8], f2[8], f3[8], f4[8];
    for (i = 0; i < 8; i++) f1[i] = float(i);
    for (i = 0; i < 8; i++) f2[i] = float(i*2);
    for (i = 0; i < 8; i++) f3[i] = float(i*3);
    for (i = 0; i < 8; i++) f4[i] = 0.0;

    MultiplyAndAdd(f1, f2, f3, f4);	

    for (i = 0; i < 8; i++) {
        printf("result[%d] = %f (%f*%f + %f)\n", i, f4[i], f1[i], f2[i], f3[i]);
    }

    return 0;
}
