#include <FastLED.h>

#define NUM_LEDS 60
#define LED_PIN 2

CRGB leds[NUM_LEDS];
size_t buf_length = NUM_LEDS * 3;

void setup()
{
	FastLED.addLeds<NEOPIXEL, LED_PIN>(leds, NUM_LEDS);

	Serial.begin(250000);
}

void loop()
{
	Serial.readBytes((char *)leds, buf_length);
	FastLED.show();
}
