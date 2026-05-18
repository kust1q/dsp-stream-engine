#include "dsp.h"
#include <cstdlib>

class BiquadFilter {
public:
    BiquadFilter(float b0, float b1, float b2, float a1, float a2)
        : b0(b0), b1(b1), b2(b2), a1(a1), a2(a2), x1(0), x2(0), y1(0), y2(0) {}

    void Process(float* buffer, int length) {
        if (!buffer || length <= 0) return;

        for (int i = 0; i < length; ++i) {
            float x = buffer[i];
            float y = b0 * x + b1 * x1 + b2 * x2 - a1 * y1 - a2 * y2;

            // Update states
            x2 = x1;
            x1 = x;
            y2 = y1;
            y1 = y;

            buffer[i] = y;
        }
    }

private:
    float b0, b1, b2, a1, a2;
    float x1, x2, y1, y2;
};

extern "C" {

BiquadHandle Biquad_Create(float b0, float b1, float b2, float a1, float a2) {
    return new BiquadFilter(b0, b1, b2, a1, a2);
}

void Biquad_Process(BiquadHandle handle, float* buffer, int length) {
    if (handle) {
        static_cast<BiquadFilter*>(handle)->Process(buffer, length);
    }
}

void Biquad_Destroy(BiquadHandle handle) {
    if (handle) {
        delete static_cast<BiquadFilter*>(handle);
    }
}

}
