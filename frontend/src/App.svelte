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
        {#if err}
          <p>{errMsg}</p>
        {:else if resultText !== ""}
          <div class="progress-card">
            <!-- Status Header -->
            <div class="card-header">
              <h3 class="status-title">{resultText}</h3>
            </div>

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
          </div>
        {/if}
      </div>
    </div>
    <div class="">
      <button class="submit-btn" on:click={() => (isModalOpen = true)}
        >Add Cookies</button
      >
    </div>
    <div class="container">
      <button class="input-btn" on:click={handlePickFolder}>
        Choose Download Location
      </button>
      <Modal bind:open={isModalOpen}>
        <h2>Add cookies</h2>
        <div class="input-div">
          <input
            type="text"
            class="input-box"
            placeholder="Paste your cookies here!"
            bind:value={cookies}
          />
        </div>
      </Modal>
    </div>
  </div>
</main>

<style>
  /* Root Layout & Background */
  .app-container {
    height: 100dvh;
    width: 100vw;
    background-color: #000;
    background-image: url("https://pixiv.re/99817334.jpg"); /* Fallback image endpoint */
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

  .input-wrapper {
    position: relative;
    width: min(64rem, 90%);
  }

  .input-field {
    width: 100%;
    padding: 1rem 7rem 1rem 1rem;
    border-radius: 9999px;
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
    border-color: #c084fc;
  }

  .input-field:focus {
    border-color: #d8b4fe;
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
    background-color: #9333ea;
    color: #fff;
    border: none;
    border-radius: 9999px;
    cursor: pointer;
    transition:
      background-color 0.2s ease,
      opacity 0.2s ease;
  }

  .submit-btn:hover:not(.disabled) {
    background-color: #a855f7;
  }

  /* Header & Navigation Bar */
  .top-nav {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    padding: 1rem;
    box-sizing: border-box;
    background: linear-gradient(
      to bottom,
      rgba(0, 0, 0, 0.75),
      rgba(0, 0, 0, 0.5),
      transparent
    );
  }
/* Glassmorphic Container Box */
  .progress-card {
    width: min(40rem, 90%);
    padding: 1.75rem 2rem;
    background-color: rgba(38, 38, 38, 0.75);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 1px solid rgba(192, 132, 252, 0.35);
    border-radius: 1.5rem;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6),
                0 0 30px rgba(147, 51, 234, 0.15);
    color: #fff;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    box-sizing: border-box;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
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
    border-radius: 9999px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.8);
  }

  /* Animated Fill Bar */
  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #7e22ce 0%, #a855f7 50%, #c084fc 100%);
    border-radius: 9999px;
    position: relative;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 0 12px rgba(168, 85, 247, 0.8);
  }

  /* Glowing Lead Edge */
  .glow-head {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 12px;
    background: #ffffff;
    border-radius: 9999px;
    filter: blur(2px);
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
</style>