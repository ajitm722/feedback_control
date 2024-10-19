import pandas as pd
import matplotlib.pyplot as plt
import matplotlib.animation as animation

# Load the data
df = pd.read_csv('temperature_data.csv')

# Parameters
t_end = 100.0
dt = 0.04
frame_amount = int(t_end / dt)

# Initialize figure and axis
fig, ax = plt.subplots(figsize=(16, 9), dpi=120)

# Plot settings
ax.set_xlim(0, t_end)
ax.set_ylim(0, 100)
ax.set_xlabel('Time')
ax.set_ylabel('Cooling System Actuator')
ax.grid(True)

# Plot objects
line_rm_1, = ax.plot([], [], 'black', linewidth=4, label='Temperature Actuation')
line_vol_r1_line, = ax.plot([], [], 'r', linewidth=2)
annotation = ax.text(0.5, 0.95, '', transform=ax.transAxes, ha='center', fontsize=12, bbox=dict(boxstyle='round', facecolor='white', edgecolor='black'))

# Update function for animation
def update_plot(num):
    if num >= len(df):
        num = len(df) - 1

    ax.set_title(f'Time: {df.iloc[num]["Time"]:.2f} seconds')

    line_rm_1.set_data(df['Time'][:num], df['Actual Temperature'][:num])
    line_vol_r1_line.set_data([0, t_end], [df['Reference Temperature'][num]] * 2)

    # Update annotation
    if pd.notna(df.iloc[num]['Annotation']) and df.iloc[num]['Annotation']:
        annotation.set_text(df.iloc[num]['Annotation'])
    else:
        annotation.set_text(f'People In Room: {int(df.iloc[num]["People In Room"])}')

    return line_rm_1, line_vol_r1_line, annotation

# Create animation
ani = animation.FuncAnimation(fig, update_plot, frames=frame_amount, interval=20, repeat=True, blit=True)

plt.show()

# MQTT (Message Queuing Telemetry Transport) is a lightweight messaging protocol designed for efficient communication in situations with limited bandwidth or unreliable networks. Here's a detailed yet straightforward explanation of MQTT internals:

# Key Components of MQTT
# Broker:

# Role: Acts as a central hub for messages. It receives messages from publishers and distributes them to subscribers.
# Function: Manages all client connections and handles message routing. Examples include Mosquitto, HiveMQ, and EMQX.
# Client:

# Role: Any device or application that publishes messages or subscribes to topics.
# Function: Can be either a publisher (sends messages) or a subscriber (receives messages) or both.
# Topic:

# Role: Acts as a message filter. Clients subscribe to topics to receive specific messages and publish messages to topics.
# Function: Topics are structured in a hierarchical manner, like folders in a file system (e.g., home/livingroom/temperature).
# Message:

# Role: The data being transmitted.
# Function: Contains the payload (actual content) and is sent to a topic.
# QoS (Quality of Service):

# QoS 0: At most once – The message is delivered once or not at all. There's no guarantee of delivery.
# QoS 1: At least once – The message is guaranteed to be delivered at least once. It might be delivered multiple times.
# QoS 2: Exactly once – The message is guaranteed to be delivered exactly once. This is the safest but also the slowest.
# How MQTT Works
# Establishing a Connection:

# A client connects to the MQTT broker using a network connection (e.g., TCP/IP). The connection process includes sending a CONNECT message and receiving a CONNACK message in response.
# Publishing Messages:

# A client sends a message to a specific topic using a PUBLISH message. The broker then forwards this message to all clients that have subscribed to that topic.
# Subscribing to Topics:

# A client subscribes to a topic to receive messages published to that topic. The subscription request is sent to the broker using a SUBSCRIBE message, and the broker responds with a SUBACK message.
# Receiving Messages:

# When a message is published to a topic, the broker checks the list of subscribers for that topic and forwards the message to each subscriber using a PUBLISH message.
# Disconnecting:

# When a client wants to disconnect, it sends a DISCONNECT message to the broker. The broker then closes the connection.
# MQTT Messages and Packet Structure
# CONNECT: Sent by a client to establish a connection to the broker.
# CONNACK: Sent by the broker to acknowledge the connection request.
# PUBLISH: Sent by a client to publish a message to a topic.
# PUBACK: Acknowledgment for QoS 1 messages.
# PUBREC: Intermediate acknowledgment for QoS 2 messages.
# PUBREL: Release acknowledgment for QoS 2 messages.
# PUBCOMP: Final acknowledgment for QoS 2 messages.
# SUBSCRIBE: Sent by a client to subscribe to a topic.
# SUBACK: Acknowledgment for a subscription request.
# UNSUBSCRIBE: Sent by a client to unsubscribe from a topic.
# UNSUBACK: Acknowledgment for an unsubscription request.
# PINGREQ: Sent by a client to keep the connection alive.
# PINGRESP: Response to the PINGREQ message.
# DISCONNECT: Sent by a client to disconnect from the broker.
# Key Features
# Lightweight: MQTT has a small header size and minimal protocol overhead, making it efficient for low-bandwidth or high-latency networks.
# Publish/Subscribe Model: Decouples the sender and receiver, allowing for more flexible and scalable communication.
# Quality of Service Levels: Provides different guarantees for message delivery, accommodating various needs for reliability and performance.
# Retained Messages: Allows the broker to store the last message sent to a topic so new subscribers can receive the most recent message immediately upon subscription.
# Last Will and Testament: Provides a way for clients to notify others if they unexpectedly disconnect.
# Use Cases
# IoT Devices: Perfect for small devices with limited resources, such as sensors and smart home devices.
# Real-time Applications: Suitable for real-time messaging needs, such as chat applications or live data feeds.
