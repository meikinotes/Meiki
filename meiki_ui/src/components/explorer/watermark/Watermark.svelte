<script context="module" lang="ts">
    export enum WatermarkType {
        INFO,
        ERROR,
    }

    function getColorClass(type: WatermarkType) {
        switch (type) {
            case WatermarkType.INFO:
                return "text-content stroke-content opacity-30"
            case WatermarkType.ERROR:
                return "text-error-content stroke-error-content opacity-30"
        }
    }
</script>

<script lang="ts">
    export let type: WatermarkType = WatermarkType.INFO
    export let message: string
</script>

<div class="flex flex-col gap-5 h-full items-center justify-center {getColorClass(type)}">
    <span class="flex-shrink-0 w-20 h-20">
        <slot />
    </span>
    <span class="px-10 font-bold text-xl text-center">
        {message}
    </span>
</div>
