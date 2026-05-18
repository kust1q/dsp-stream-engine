#ifndef DSP_H
#define DSP_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque handle for the Biquad filter
typedef void* BiquadHandle;

// Create a new Biquad filter instance
// Coefficients: b0, b1, b2, a1, a2
BiquadHandle Biquad_Create(float b0, float b1, float b2, float a1, float a2);

// Process a chunk of samples in-place
void Biquad_Process(BiquadHandle handle, float* buffer, int length);

// Destroy the instance
void Biquad_Destroy(BiquadHandle handle);

#ifdef __cplusplus
}
#endif

#endif // DSP_H
