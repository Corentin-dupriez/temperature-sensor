#include <Adafruit_Sensor.h>
#include <DHT.h>
#include <DHT_U.h>

#define DHTPIN 2

#define DHTTYPE DHT11

DHT_Unified dht(DHTPIN, DHTTYPE);

uint32_t delayMS;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  dht.begin();

  sensor_t sensor;

  dht.temperature().getSensor(&sensor);
  delayMS = sensor.min_delay / 1000;
}

void loop() {
  // put your main code here, to run repeatedly:
  delay(delayMS);
  sensors_event_t event;

  dht.temperature().getEvent(&event);

  Serial.print(F("Temp:"));
  Serial.print(event.temperature);

  dht.humidity().getEvent(&event);

  Serial.print(F(" Humidity:"));
  Serial.print(event.relative_humidity);
  Serial.println(F("%"));
}
