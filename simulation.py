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
