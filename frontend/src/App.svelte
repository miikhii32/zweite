<script lang="ts">
  import { Rip, SelectTargetDirectory, GetDefaultDocumentsFolder } from '../wailsjs/go/main/App.js'
  import { onDestroy, onMount } from 'svelte';
  import { EventsOn } from "../wailsjs/runtime"
  import Modal from '../components/Modal.svelte';

  let resultText: string = ""
  let url: string = "";
  let cookies: string = "";
  let err: boolean = false;
  let errMsg: string = "";
  let unsubError
  let unsubTitle
  let unsubPage
  let selectedFolder = "";
  let statusMessage = "";
  let isModalOpen = false;

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
    Rip(url, cookies, selectedFolder)
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
      <div class="msgs">
        {#if err}
          <p>{errMsg}</p>
        {:else}
          <p>{resultText}</p>
        {/if}
      </div>
      <div class="input-div">
        <button class="input-btn" on:click={greet}>Submit!</button>
      </div>
      <div class="input-div">
        <button class="input-btn" on:click={() => isModalOpen = true}>Add Cookies</button>
      </div>
    </div>
    <div class="container">
    <button class="input-btn" on:click={handlePickFolder}>
        Choose Download Location
    </button>
    <Modal bind:open={isModalOpen}>
      <h2>Add cookies</h2>
      <div class="input-div">
        <input type="text" class="input-box" placeholder="Paste your cookies here!" bind:value={cookies}>
      </div>
  </Modal>
  </div>
  </div>
</main>

<style>

  h1 {
    color: white;
  }
  
  h2 {
    color: white;
  }
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

  p {
    text-align: center;
  }

  .msgs {
    padding-bottom: 2rem;
  }

</style>
