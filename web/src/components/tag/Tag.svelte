<script>
    import "./tag.css"
    export let value = ""
    export let deletable = false
</script>


<span class="tag">
    {value}

    {#if deletable}
        <i class="close" title="Delete">x</i>
    {/if}
</span>