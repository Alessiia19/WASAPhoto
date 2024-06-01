<script>
import animationData from '@/assets/not_found_animation.json'

export default {
    data: function () {
        return {
            animationData,
        }
    },

}
</script>

<template>
    <!-- Header -->
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
        <h1 class="app-title unselectable">WASAPhoto</h1>
    </header>

    <div class="not-found-container">
        <lottie-player :src="animationData" background="transparent" speed="0.5" style="width: 300px; height: 300px;"
            loop autoplay></lottie-player>
        <h1>404 Not Found</h1>
        <p style="font-size: 18px;">The page you are looking for doesn't exist or has been moved.</p>
    </div>
</template>

<style scoped>
.not-found-container {
    width: calc(100% + 230px);
    height: 100%;
    margin-top: 170px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    font-family: 'Fredoka', sans-serif;
    font-weight: 400;
}
</style>