<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import { Rip, SelectTargetDirectory, GetDefaultDocumentsFolder } from '../wailsjs/go/main/App.js'
  import { onDestroy, onMount } from 'svelte';
  import { EventsOn } from "../wailsjs/runtime"

  let resultText: string = ""
  let url: string = "";
  let err: boolean = false;
  let errMsg: string = "";
  let unsubError
  let unsubTitle
  let unsubPage
  let selectedFolder = "";
  let statusMessage = "";

  onMount(() => {
    unsubError = EventsOn("error-emit", (data) => {
      err = true;
      errMsg = data
    })
    unsubTitle = EventsOn("title-get", (data) => {
      err = false;
      resultText = data
    })
    // unsubPage = EventsOn("", () => {})
  })

  onDestroy(() => {
    if (unsubError) unsubError();
    if (unsubTitle) unsubTitle();
  })

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
    if (selectedFolder === ""){
      selectedFolder = await GetDefaultDocumentsFolder()
    }
    Rip(url, "", selectedFolder)
  }
</script>

<main>
  <div>
    <h1>New Hakuneko!</h1>
    <p>Hello there! Paste the link to a chapter you'd like to download!</p>
    <div>
      <div class="input-div">
        <input type="text" class="input-box" placeholder="Paste your link here!" bind:value={url}>
      </div>
      <div class="input-div">
        <button class="input-btn" on:click={greet}>Submit!</button>
      </div>
      {#if err}
        <p>{errMsg}</p>
      {:else}
        <p>{resultText}</p>
      {/if}
    </div>
    <div class="container">
    <button on:click={handlePickFolder}>
        Choose Download Location
    </button>
  </div>
  </div>
</main>

<style>

  .input-div {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 2rem;
  }

  .input-btn{
    width: 7rem;
    height: 3rem;
  }

  .input-box{
      background: solid;
      background-color: #222E42;
      width: 30rem;
      height: 2rem;
      text-align: center;
      color: white;
      font-size: 0.95rem;
      border: none;

  }
  
  .input-box:focus{
    background: none;
    outline: none;
  }
  
  .input-box:hover{
    background: none;
    outline: none;
  }

</style>
