# Mystreaming

## Project Overview  
Before presenting the project, let's understand **HLS** and its advantages.  

The **HLS (HTTP Live Streaming)** protocol delivers video over **HTTP** by dividing it into small segments and using an **`.m3u8`** playlist for continuous playback. Its main advantage is **adaptive bitrate (ABR) streaming**, which adjusts the quality based on the user's connection.  

### Features  

- **Video Processing with FFmpeg**  
  - Uses **FFmpeg** to process videos for the **HLS protocol**.  
  - Displays **processing progress** and provides an estimated **completion time**.  
  - In case of an error, allows **resuming FFmpeg** from where it left off.  

- **Chunked File Upload**  
  - Enables **segmented video uploads**, allowing users to resume the upload **at any time** if interrupted.  

- **Authentication & Database**  
  - Implements **administrator and user authentication**.  
  - Uses **PostgreSQL** as the database.  
  - Secures authentication with **JWT tokens**.  
