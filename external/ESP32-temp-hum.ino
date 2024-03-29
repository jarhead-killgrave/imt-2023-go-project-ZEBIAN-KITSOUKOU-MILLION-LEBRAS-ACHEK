#include "DHT.h"
#include <WiFi.h>
#include <ArduinoJson.h>
#include <time.h>
#include <ctime>
#include <cstdlib>
extern "C" {
  #include "freertos/FreeRTOS.h"
  #include "freertos/timers.h"
}
#include <AsyncMqttClient.h>

#define WIFI_SSID "Enter your SSID"
#define WIFI_PASSWORD "Enter your WiFi password"

// Raspberry Pi Mosquitto MQTT Broker
#define MQTT_HOST IPAddress(172, 20, 1, xxx)
#define MQTT_PORT 1883

// Temperature MQTT Topics
#define MQTT_PUB_TEMP "airport/T"
#define MQTT_PUB_HUM  "airport/H"

// Digital pin connected to the DHT sensor
#define DHTPIN 4  

#define DHTTYPE DHT11

#define ID "Enter your ID"
#define PASS "Enter your password"

char buffer[80];
time_t now = time(nullptr);

enum class MeasurementType {
  TEMPERATURE,
  HUMIDITY
  // Add other types as needed
};

struct Measurement {
  int64_t sensor_id;
  String airport_id;
  MeasurementType type;
  double value;
  String unit;
  time_t timestamp;  // This is the C++ equivalent of Go's time.Time
};

Measurement measureH;
Measurement measureT;

String jsonH;
String jsonT;

// Creation of Json object
StaticJsonDocument<200> toJsonH;
StaticJsonDocument<200> toJsonT;

// Initialize DHT sensor
DHT dht(DHTPIN, DHTTYPE);

// Variables to hold sensor readings
float temp;
float hum;
int QoS = 1;

AsyncMqttClient mqttClient;
TimerHandle_t mqttReconnectTimer;
TimerHandle_t wifiReconnectTimer;

unsigned long previousMillis = 0;   // Stores last time temperature was published
const long interval = 10000;        // Interval at which to publish sensor readings

time_t current_time = time(nullptr);

const char* ntpServer = "pool.ntp.org";
const long  gmtOffset_sec = 3600;
const int   daylightOffset_sec = 3600;

struct tm timeinfo;

void printLocalTime()
{
  //struct tm timeinfo;
  if(!getLocalTime(&timeinfo)){
    Serial.println("Failed to obtain time");
    return;
  }
  Serial.println(&timeinfo, "%A, %B %d %Y %H:%M:%S");
}

void connectToWifi() {
  Serial.println("Connecting to Wi-Fi...");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
}

void connectToMqtt() {
  Serial.println("Connecting to MQTT...");
  mqttClient.setCredentials(ID, PASS);
  mqttClient.connect();
}

void WiFiEvent(WiFiEvent_t event) {
  Serial.printf("[WiFi-event] event: %d\n", event);
  switch(event) {
    case SYSTEM_EVENT_STA_GOT_IP:
      Serial.println("WiFi connected");
      Serial.println("IP address: ");
      Serial.println(WiFi.localIP());
      connectToMqtt();
      configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
      printLocalTime();
      break;
    case SYSTEM_EVENT_STA_DISCONNECTED:
      Serial.println("WiFi lost connection");
      xTimerStop(mqttReconnectTimer, 0); // ensure we don't reconnect to MQTT while reconnecting to Wi-Fi
      xTimerStart(wifiReconnectTimer, 0);
      break;
  }
}

void onMqttConnect(bool sessionPresent) {
  Serial.println("Connected to MQTT.");
  Serial.print("Session present: ");
  Serial.println(sessionPresent);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {
  Serial.println("Disconnected from MQTT.");
  if (WiFi.isConnected()) {
    xTimerStart(mqttReconnectTimer, 0);
  }
}

void onMqttPublish(uint16_t packetId) {
  Serial.print("Publish acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
}

void setup() {
  Serial.begin(115200);
  Serial.println();

  dht.begin();

  measureH.sensor_id = 1;
  measureH.airport_id = "CDG"; 
  measureH.type = MeasurementType::TEMPERATURE; 
  measureH.value = 25.0; 
  measureH.unit = "C";  
  measureH.timestamp = time(nullptr); 

  measureT.sensor_id = 1;
  measureT.airport_id = "CDG"; 
  measureT.type = MeasurementType::HUMIDITY; 
  measureT.value = 40.0; 
  measureT.unit = "%";  
  measureT.timestamp = time(nullptr); 


  toJsonH["sensor_id"] = measureH.sensor_id;
  toJsonH["airport_id"] = measureH.airport_id.c_str();
  toJsonH["type"] = "HUMIDITY";
  toJsonH["value"] = measureH.value;
  toJsonH["unit"] = measureH.unit.c_str();
  toJsonH["timestamp"] = measureH.timestamp;

  toJsonT["sensor_id"] = measureT.sensor_id;
  toJsonT["airport_id"] = measureT.airport_id.c_str();
  toJsonT["type"] = "TEMPERATURE";
  toJsonT["value"] = measureT.value;
  toJsonT["unit"] = measureT.unit.c_str();
  toJsonT["timestamp"] = measureT.timestamp;

  // Init and get the time
  configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
  printLocalTime();

  dht.begin();
  
  mqttReconnectTimer = xTimerCreate("mqttTimer", pdMS_TO_TICKS(2000), pdFALSE, (void*)0, reinterpret_cast<TimerCallbackFunction_t>(connectToMqtt));
  wifiReconnectTimer = xTimerCreate("wifiTimer", pdMS_TO_TICKS(2000), pdFALSE, (void*)0, reinterpret_cast<TimerCallbackFunction_t>(connectToWifi));

  WiFi.onEvent(WiFiEvent);

  mqttClient.onConnect(onMqttConnect);
  mqttClient.onDisconnect(onMqttDisconnect);
  mqttClient.onPublish(onMqttPublish);
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
  mqttClient.setCredentials(ID, PASS);
  connectToWifi();
}

void loop() {
  unsigned long currentMillis = millis();
  if (currentMillis - previousMillis >= interval) {
    // Save the last time a new reading was published
    previousMillis = currentMillis;

    // Read temperature as Celsius
    temp = dht.readTemperature();

    // New DHT sensor readings
    hum = dht.readHumidity();

    toJsonH["value"] = hum;
    toJsonT["value"] = temp;

    strftime(buffer, sizeof(buffer), "%Y-%m-%dT%H:%M:%S", &timeinfo);
    String timestamp = String(buffer);
    toJsonH["timestamp"] = timestamp;
    toJsonT["timestamp"] = timestamp;
    //temp = dht.readTemperature(true); If Farhenheit is needed


    serializeJson(toJsonH, jsonH);
    serializeJson(toJsonT, jsonT);

    // Check if any reads failed and exit early (to try again).
    if (isnan(temp) || isnan(hum)) {
      Serial.println(F("Failed to read from DHT sensor!"));
      return;
    }
    
    // Publish an MQTT message on topic esp32/dht/temperature
    uint16_t packetIdPub1 = mqttClient.publish(MQTT_PUB_TEMP, QoS, true, jsonT.c_str());                            
    Serial.printf("Publishing on topic %s at QoS 1, packetId: %i", MQTT_PUB_TEMP, packetIdPub1);
    Serial.printf("Message: %.2f \n", temp);

    // Publish an MQTT message on topic esp32/dht/humidity
    uint16_t packetIdPub2 = mqttClient.publish(MQTT_PUB_HUM, QoS, true, jsonH.c_str());                            
    Serial.printf("Publishing on topic %s at QoS 1, packetId %i: ", MQTT_PUB_HUM, packetIdPub2);
    Serial.printf("Message: %.2f \n", hum);

    // Init and get the time
    configTime(gmtOffset_sec, daylightOffset_sec, ntpServer);
    printLocalTime();
  }
}