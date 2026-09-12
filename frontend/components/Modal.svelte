<script>
  // Expose an 'open' prop to control visibility from the parent
  export let open = false;

  let dialogElement;

  // Reactively open/close the native dialog based on the prop
  $: if (dialogElement) {
    if (open) {
      dialogElement.showModal(); // Opens the dialog with a native backdrop
    } else {
      dialogElement.close();
    }
  }
</script>

<dialog bind:this={dialogElement} on:close={() => (open = false)} class="modal">
  <div class="modal-content">
  <h1 style="color: white;">Add cookies</h1>
    <div class="modal-slot">
      <slot />
  
    </div>
    <div class="modal-btn">
      <button class="other-btn" on:click={() => (open = false)}> Close </button>
    </div>
  </div>
</dialog>
<style>
  /* Styles the modal dialog window itself */
  .modal {
    border: none;
    background-color: #222E42;
    border-radius: 8px;
    padding: 0;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    width: 90%;
    height: 90%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  /* Trigger the opening animation when the native [open] attribute is applied */
  .modal[open] {
    animation: scaleUp 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  }

  /* Style the native backdrop overlay */
  .modal::backdrop {
    background-color: rgba(27, 38, 54, 1);
    backdrop-filter: blur(4px);
  }

  /* Trigger the backdrop fade-in when open */
  .modal[open]::backdrop {
    animation: fadeIn 0.25s ease-out forwards;
  }

  .modal-content{
    width: 100%;
  }

  .modal-slot {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .modal-btn {
    display: flex;
    justify-content: center;
    align-items: center;
  }


  .other-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 16px;
    padding: 8px 16px;
    font-weight: 600;
    font-size: small;
    background-color: #1e1b4b;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    width: 15rem;
    height: 2.5rem;
    transition:
      background-color 0.2s ease,
      opacity 0.2s ease; 
  }

  .other-btn:hover:not(.disabled) {
    background-color: #3b82f6;
  }

  /* --- Keyframe Animations --- */
  @keyframes scaleUp {
    from {
      transform: scale(0.9) translateY(20px);
      opacity: 0;
    }
    to {
      transform: scale(1) translateY(0);
      opacity: 1;
    }
  }

  

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
</style>