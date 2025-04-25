package ginservice

const apiDoc = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" rel="stylesheet">
</head>

<body class="bg-gray-100 font-sans">
    <div class="container mx-auto p-8">
        <h1 class="text-3xl font-bold text-blue-600 mb-6">{{.title}}</h1>
        <h2 class="text-2xl font-bold text-gray-800 mb-4">目录</h2>
        <ul class="list-disc pl-6 mb-6">
            {{range $index, $api :=.apis}}
            <li><a href="#api-{{$index}}" class="text-blue-600 hover:underline">{{$api.Method}} {{$api.URL}}</a></li>
            {{end}}
        </ul>
        {{range $index, $api :=.apis}}
        <div id="api-{{$index}}" class="bg-white p-6 rounded shadow mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">{{$api.Method}} {{$api.URL}}</h2>
            <p class="text-gray-600 mb-4">Handler: {{$api.HandlerName}}</p>
            {{if $api.Request}}
            <h3 class="text-xl font-bold text-gray-700 mb-3">Request</h3>
            <p class="text-gray-600 mb-2">Name: {{$api.Request.Name}}</p>
            <p class="text-gray-600 mb-2">PkgPath: {{$api.Request.PkgPath}}</p>
            <p class="text-gray-600 mb-4">Curl: <code>{{$api.Request.CurlString}}</code></p>
            <table class="table-auto w-full mb-6">
                <thead>
                    <tr>
                        <th class="px-4 py-2 border">Field Name</th>
                        <th class="px-4 py-2 border">Type</th>
                        <th class="px-4 py-2 border">Required</th>
                        <th class="px-4 py-2 border">Tag</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $api.Request.FieldInfos}}
                    <tr>
                        <td class="px-4 py-2 border">{{.Name}}</td>
                        <td class="px-4 py-2 border">{{.Typ}}</td>
                        <td class="px-4 py-2 border">{{.Required}}</td>
                        <td class="px-4 py-2 border">{{.Tag}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <h4 class="text-lg font-bold text-gray-700 mb-2">Request JSON 样例</h4>
            <pre class="bg-gray-200 p-4 rounded mb-6" id="request-json-sample-{{$index}}"></pre>
            <script>
                const curlString = '{{$api.Request.CurlString}}';
                const jsonStart = curlString.indexOf('{');
                const jsonEnd = curlString.lastIndexOf('}');
                if (jsonStart!== -1 && jsonEnd!== -1) {
                    const json = curlString.substring(jsonStart, jsonEnd + 1);
                    try {
                        const parsedJson = JSON.parse(json);
                        const formattedJson = JSON.stringify(parsedJson, null, 4);
                        document.getElementById('request-json-sample-{{$index}}').textContent = formattedJson;
                    } catch (error) {
                        document.getElementById('request-json-sample-{{$index}}').textContent = json;
                    }
                }
            </script>
            {{else}}
            <p class="text-gray-600 mb-4">无请求参数</p>
            {{end}}
            {{if $api.Response}}
            <h3 class="text-xl font-bold text-gray-700 mb-3">Response</h3>
            <table class="table-auto w-full mb-6">
                <thead>
                    <tr>
                        <th class="px-4 py-2 border">Field Name</th>
                        <th class="px-4 py-2 border">Type</th>
                        <th class="px-4 py-2 border">Tag</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $api.Response.FieldInfos}}
                    <tr>
                        <td class="px-4 py-2 border">{{.Name}}</td>
                        <td class="px-4 py-2 border">{{.Typ}}</td>
                        <td class="px-4 py-2 border">{{.Tag}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            <h4 class="text-lg font-bold text-gray-700 mb-2">Response JSON 样例</h4>
            <pre class="bg-gray-200 p-4 rounded mb-6" id="response-json-sample-{{$index}}"></pre>
            
            {{else}}
            <p class="text-gray-600 mb-4">无响应参数</p>
            {{end}}
        </div>
        {{end}}
    </div>
</body>

</html>
`

//const apiDoc = `<!DOCTYPE html>
//<html lang="en">
//
//<head>
//
//	<meta charset="UTF-8">
//	<meta name="viewport" content="width=device-width, initial-scale=1.0">
//	<title>API Doc</title>
//	<script src="https://cdn.tailwindcss.com"></script>
//	<link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css" rel="stylesheet">
//	<style>
//	    .method {
//	        @apply text-white px-4 py-2 mr-2 rounded-md font-bold min-w-[80px] inline-block text-center;
//	    }
//
//	    .GET {
//	        @apply bg-blue-500;
//	    }
//
//	    .POST {
//	        @apply bg-green-500;
//	    }
//
//	    .PUT {
//	        @apply bg-yellow-500;
//	    }
//
//	    .DELETE {
//	        @apply bg-red-500;
//	    }
//
//	    /* 动画效果 */
//	    .menu li {
//	        transition: all 0.3s ease;
//	    }
//
//	    .menu li:hover {
//	        transform: scale(1.05);
//	        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
//	    }
//
//	    .tab-pane {
//	        animation: fadeIn 0.5s ease;
//	    }
//
//	    @keyframes fadeIn {
//	        from {
//	            opacity: 0;
//	        }
//
//	        to {
//	            opacity: 1;
//	        }
//	    }
//	</style>
//
//</head>
//
//<body class="bg-gray-100">
//
//	<nav class="bg-gray-900 fixed top-0 w-full z-10 shadow-md">
//	    <div class="container mx-auto px-4 py-3 flex justify-between items-center">
//	        <a href="#" class="text-white text-xl font-bold">API Documentation</a>
//	        <button id="menu-toggle" class="text-white md:hidden focus:outline-none">
//	            <i class="fas fa-bars"></i>
//	        </button>
//	    </div>
//	</nav>
//	<div class="container mx-auto px-4 py-20 grid grid-cols-12 gap-4 md:gap-8">
//	    <div class="col-span-12 md:col-span-4 menu">
//	        <ul class="space-y-2">
//	            <li class="bg-white rounded-md shadow-md p-3 flex items-center cursor-pointer" data-target="api1">
//	                <span class="method GET">GET</span>
//	                <span class="ml-2 text-gray-700 font-bold">/api/users</span>
//	            </li>
//	            <li class="bg-white rounded-md shadow-md p-3 flex items-center cursor-pointer" data-target="api2">
//	                <span class="method POST">POST</span>
//	                <span class="ml-2 text-gray-700 font-bold">/api/users</span>
//	            </li>
//	        </ul>
//	    </div>
//	    <div class="col-span-12 md:col-span-8">
//	        <div id="api1" class="bg-white rounded-md shadow-md p-6 tab-pane hidden">
//	            <h4 class="text-2xl font-bold mb-4 text-gray-800">Basic Info</h4>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Key</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Value</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Handler Name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">GetUsersHandler</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Method</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">GET</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Route</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">/api/users</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Return Type</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">List[User]</td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	            <h4 class="text-2xl font-bold mb-4 mt-6 text-gray-800">Request Info</h4>
//	            <p class="text-gray-600">Not Defined</p>
//	            <h4 class="text-2xl font-bold mb-4 mt-6 text-gray-800">Response Fields</h4>
//	            <p class="text-gray-600">Description: Returns a list of users</p>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Name</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Type</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Tag</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">id</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">int</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">required</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">str</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">required</td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	        </div>
//	        <div id="api2" class="bg-white rounded-md shadow-md p-6 tab-pane hidden">
//	            <h4 class="text-2xl font-bold mb-4 text-gray-800">Basic Info</h4>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Key</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Value</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Handler Name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">CreateUserHandler</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Method</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">POST</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Route</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">/api/users</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Return Type</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">User</td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	            <h4 class="text-2xl font-bold mb-4 mt-6 text-gray-800">Request Info</h4>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Key</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Value</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Request Name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">CreateUserRequest</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Package Path</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">com.example.api</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">Example Url</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600"><pre>curl -X POST /api/users -d '{"name": "John Doe"}'</pre></td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	            <h5 class="text-lg font-bold mb-2 mt-4 text-gray-800">Request Fields</h5>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Name</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Type</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Tag</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Required</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">str</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">required</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">true</td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	            <h4 class="text-2xl font-bold mb-4 mt-6 text-gray-800">Response Fields</h4>
//	            <p class="text-gray-600">Description: Returns the created user</p>
//	            <table class="table-auto w-full border border-gray-300">
//	                <thead>
//	                    <tr>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Name</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Type</th>
//	                        <th class="border border-gray-300 p-3 text-left text-gray-700">Tag</th>
//	                    </tr>
//	                </thead>
//	                <tbody>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">id</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">int</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">required</td>
//	                    </tr>
//	                    <tr>
//	                        <td class="border border-gray-300 p-3 text-gray-600">name</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">str</td>
//	                        <td class="border border-gray-300 p-3 text-gray-600">required</td>
//	                    </tr>
//	                </tbody>
//	            </table>
//	        </div>
//	    </div>
//	</div>
//	<script>
//	    const menuItems = document.querySelectorAll('.menu li');
//	    const tabPanes = document.querySelectorAll('.tab-pane');
//
//	    menuItems.forEach(item => {
//	        item.addEventListener('click', () => {
//	            const targetId = item.dataset.target;
//	            tabPanes.forEach(pane => {
//	                if (pane.id === targetId) {
//	                    pane.classList.remove('hidden');
//	                } else {
//	                    pane.classList.add('hidden');
//	                }
//	            });
//	        });
//	    });
//
//	    const menuToggle = document.getElementById('menu-toggle');
//	    const menu = document.querySelector('.menu');
//	    menuToggle.addEventListener('click', () => {
//	        menu.classList.toggle('hidden');
//	    });
//	</script>
//
//</body>
//
//</html>`
