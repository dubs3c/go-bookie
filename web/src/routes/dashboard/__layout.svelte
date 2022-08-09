<script>
    import "../../static/css/grid.css"
    import "../../static/css/global.css"
</script>
<div class="container">
<main>
    
        <div class="row hr">
            <div class="col">
                <h1><a href="/" style="text-decoration:none; color:turquoise;">SAMLA</a></h1>
            </div>
            <nav>
                <div class="col">
                    <ul>
                        <li><a href="\settings">Settings</a></li>
                        <li><a href="\logout" class="logout-red">Logout</a></li>
                    </ul> 
                </div>
            </nav>
        </div>
    <br />

    <slot></slot>

</main>
</div>