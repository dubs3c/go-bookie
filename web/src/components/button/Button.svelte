<script>
    import './button.css';

    export let value = "Submit"
    export let large = false
    export let fullWidth = false
    export let danger = false
</script>

<button
    class="btn"
    class:large={large}
    class:fullWidth={fullWidth}
    class:danger={danger}>
    
    {value}</button>