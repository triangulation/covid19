<script>
  import Chart from "./Chart.svelte"
</script>

<main>
  <h1>COVID 19 stats</h1>
  <Chart />
</main>

<style>
  main {
    margin: 0 auto;
    padding: 0.5em;
    text-align: center;
  }
  h1 {
    font-size: 2.25rem;
    font-weight: 400;
    padding: 0;
    margin: 0;
  }

  @media (min-width: 40em) {
    main {
      padding: 1rem;
    }
    h1 {
      font-size: 3rem;
    }
  }
</style>
