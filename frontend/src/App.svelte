<script lang="ts">
  import {
    Rip,
    SelectTargetDirectory,
    GetDefaultDocumentsFolder,
  } from "../wailsjs/go/main/App.js";
  import { onDestroy, onMount } from "svelte";
  import { EventsOn } from "../wailsjs/runtime";
  import Modal from "../components/Modal.svelte";

  let resultText: string = "";
  let url: string = "";
  let cookies: string = "";
  let err: boolean = false;
  let errMsg: string = "";
  let unsubError;
  let unsubTitle;
  let unsubPage;
  let unsubTotalPages;
  let selectedFolder = "";
  let statusMessage = "";
  let isModalOpen = false;
  let downloadedPages = 0;
  let totalPages;

  onMount(() => {
    unsubError = EventsOn("error-emit", (data) => {
      err = true;
      errMsg = data;
    });
    unsubTitle = EventsOn("title-get", (data) => {
      err = false;
      resultText = data;
    });
    unsubPage = EventsOn("new-page", (data) => {
      err = false;
      downloadedPages = downloadedPages + 1;
    });
    unsubTotalPages = EventsOn("total-page", (data) => {
      totalPages = data;
    });
    // unsubPage = EventsOn("", () => {})
  });

  onDestroy(() => {
    if (unsubError) unsubError();
    if (unsubTitle) unsubTitle();
  });

  async function handlePickFolder() {
    try {
      const folderPath = await SelectTargetDirectory();
      if (folderPath) {
        selectedFolder = folderPath;
        statusMessage = "Folder selected successfully!";
      } else {
        statusMessage = "Folder selection cancelled.";
      }
    } catch (err) {
      statusMessage = "Error picking folder: " + err;
    }
  }

  async function greet(): Promise<void> {
    if (selectedFolder === "") {
      selectedFolder = await GetDefaultDocumentsFolder();
    }
    resultText = "";
    downloadedPages = 0;
    totalPages = 0;
    resultText = "Finding your raws... give us a second";
    Rip(url, cookies, selectedFolder);
    url = "";
  }
</script>

<main>
  <div class="app-container">
    <div class="top-nav">
      <h1>zwite</h1>
    </div>
    <div class="center-content">
      <div class="main-card-stack">
        <div class="input-wrapper">
          <input
            type="text"
            class="input-field"
            placeholder="Paste your link here!"
            bind:value={url}
          />
          <div class="button-wrapper">
            <button class="submit-btn" on:click={greet}>Submit!</button>
          </div>
        </div>
        <div>
          {#if err}
            <p>{errMsg}</p>
          {:else if resultText !== ""}
            <div class="progress-card">
              <!-- Status Header -->
              <div class="card-header">
                <h3 class="status-title">{resultText}</h3>
              </div>
              {#if resultText !== "Finding your raws... give us a second"}
                <!-- Progress Track & Fill -->
                <div class="track-wrapper">
                  <div
                    class="progress-track"
                    role="progressbar"
                    aria-valuenow={(downloadedPages / totalPages) * 100}
                    aria-valuemin="0"
                    aria-valuemax="100"
                  >
                  <div
                    class="progress-fill"
                    style="width: {(downloadedPages / totalPages) * 100}%"
                  >
                    <div class="glow-head"></div>
                  </div>
                </div>
              </div>
              <div class="card-footer">
                <span class="percentage">{downloadedPages}/{totalPages}</span>
              </div>
              {/if}
            </div>
          {/if}
      </div>
          <div>
            <div class="other-btn-containers">
              <button class="other-btn" on:click={() => (isModalOpen = true)}>Add Cookies</button>
              <button class="other-btn" on:click={handlePickFolder}>
                Choose Download Location
              </button>
            </div>
            {#if isModalOpen}
            <Modal bind:open={isModalOpen}>
                  <textarea
                    class="input-field-modal"
                    placeholder="Paste your cookies here!"
                    bind:value={cookies}
                  />
              </Modal>
            {/if}
          </div>
      </div>
    </div>
  </div>
</main>

<style>
  /* Root Layout & Background */
  .app-container {
    height: 100dvh;
    width: 100vw;
    background: linear-gradient(
      to bottom,
      rgba(15, 23, 42, 0.75),
      rgba(30, 41, 59, 0.35),
      transparent
    );
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
    color: #fff;
    font-family:
      system-ui,
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      sans-serif;
    position: relative;
    overflow: hidden;
  }

  /* Center Content Area */
  .center-content {
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .main-card-stack {
    width: min(84rem, 90%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem; /* Controls space between input box and progress card */
  }

  .input-wrapper {
    position: relative;
    width: min(64rem, 90%);
  }

  .input-field {
    width: 100%;
    padding: 1rem 7rem 1rem 1rem;
    border-radius: 4px;
    background-color: rgba(82, 82, 82, 0.9);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 2px solid transparent;
    color: #fff;
    outline: none;
    font-size: 1rem;
    box-sizing: border-box;
    transition: border-color 0.2s ease;
  }

  .input-field::placeholder {
    color: #d4d4d4;
  }

  .input-field:hover {
    border-color: #3b82f6;
  }

  .input-field:focus {
    border-color: #93c5fd;
  }
 
  .input-field-modal {
    width: 64%;
    height: 7rem;
    padding: 1rem 7rem 1rem 1rem;
    border-radius: 4px;
    background-color: rgba(82, 82, 82, 0.9);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 2px solid transparent;
    color: #fff;
    outline: none;
    font-size: 1rem;
    box-sizing: border-box;
    transition: border-color 0.2s ease;
  }

  .input-field-modal::placeholder {
    color: #d4d4d4;
  }

  .input-field-modal:hover {
    border-color: #3b82f6;
  }

  .input-field-modal:focus {
    border-color: #93c5fd;
  }

  .button-wrapper {
    position: absolute;
    right: 0.5rem;
    top: 0;
    bottom: 0;
    display: flex;
    align-items: center;
  }

  .submit-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    font-weight: 600;
    background-color: #1e1b4b;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition:
      background-color 0.2s ease,
      opacity 0.2s ease;
  }

  .submit-btn:hover:not(.disabled) {
    background-color: #3b82f6;
  }
  .other-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    font-weight: 600;
    font-size: small;
    background-color: #1e1b4b;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    width: 15rem;
    height: 2.5rem;
    text-align: center;
    transition:
      background-color 0.2s ease,
      opacity 0.2s ease; 
  }

  .other-btn:hover:not(.disabled) {
    background-color: #3b82f6;
  }

  /* Header & Navigation Bar */
  .top-nav {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    padding: 1rem;
    box-sizing: border-box; 
    justify-content: center;
  }
/* Glassmorphic Container Box */
  .progress-card {
    width: 50rem;
    padding: 1.75rem 2rem;
    background-color: #1e1b4b;
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 1px solid rgba(99, 102, 241, 0.35);
    border-radius: 1.5rem;
    box-shadow: 
  0 20px 40px rgba(0, 0, 0, 0.6),
  0 0 30px rgba(30, 27, 75, 0.4);
    color: #fff;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    box-sizing: border-box;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    border-radius: 4px;
    margin-bottom: 2rem;
    margin-top: 1rem;
  }

  /* Header Section */
  .card-header {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    text-align: center;
  }

  .status-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 700;
    letter-spacing: -0.01em;
    color: #f5f5f5;
  }

  /* Outer Track */
  .track-wrapper {
    width: 100%;
  }

  .progress-track {
    width: 100%;
    height: 1.25rem;
    background-color: rgba(23, 23, 23, 0.8);
    border-radius: 4px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.8);
  }

  /* Animated Fill Bar */
  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #1d4ed8 0%, #3b82f6 50%, #93c5fd 100%);
    border-radius: 4px;
    position: relative;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* Glowing Lead Edge */
  .glow-head {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 20px;
    background: #ffffff;
    border-radius: 4px;
    opacity: 0.8;
  }

  /* Bottom Percentage Indicator */
  .card-footer {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .percentage {
    font-size: 0.95rem;
    font-weight: 600;
    color: #d8b4fe;
    letter-spacing: 0.02em;
  }

 .other-btn-containers {
    display: flex;
    flex-direction: row;
    gap: 0.75rem; /* Controls distance between the buttons */
    width: 100%;
    justify-content: center;
 }

  @media (max-width: 840px) {
    .progress-card {
      width: 30rem;
    }
  }


  
  @media (max-width: 640px) {
    .progress-card {
      width: 20rem;
    }
  }
  
  @media (max-width: 440px) {
    .progress-card {
      width: 10rem;
    }
  }
</style>