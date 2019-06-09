#include <FastLED.h>

#define NUM_LEDS 30
#define LED_PIN 2

byte buf[4] = {0, 0, 0, 0};
CRGB leds[NUM_LEDS];

void setup()
{
	FastLED.addLeds<NEOPIXEL, LED_PIN>(leds, NUM_LEDS);

	Serial.begin(250000);
}

void loop()
{
	Serial.readBytes((char *)leds, NUM_LEDS * 3);
	FastLED.show();
}
