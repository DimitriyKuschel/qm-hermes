async function updateQueueTable() {
    const messages = document.getElementById('messages');

    try {
        const messagesResponse = await fetch(`/api/v1/message/list?queue_name=${queueName}`);
        const jsonDetailedQueues = await messagesResponse.json();
        messages.innerHTML = '';
        jsonDetailedQueues.forEach((item) => {
            messages.innerHTML += `        <div class="card card-flush bgi-no-repeat bgi-size-contain bgi-position-x-end h-md-50 mb-1 mb-xl-10" style="background-color: #317f7f;">
            <!--begin::Card body-->
            <div class="card-body align-items-end">
                <!--begin::Amount-->
                <span class="fs-2hx fw-bold text-white me-2 lh-1 ls-n2" id="messageCount">                
                ${item}
                </span>
                <!--end::Amount-->
            </div>
            <!--end::Card body-->
        </div>`
        });
    } catch (error) {
        console.error(error);
    }
}

(async () => await updateQueueTable())();