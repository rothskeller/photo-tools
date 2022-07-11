<script lang="ts">
  export let id: string
  export let label: string
</script>

<label for={id}>{label}</label>

<style>
  label {
    display: block;
    margin-top: 0.75rem;
    color: #444;
  }
</style>
