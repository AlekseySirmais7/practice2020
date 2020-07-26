import FetchModule from "./network";
import CurrentUser from "./currentUser";

function showModalCreatePost() {
    const darkLayer = document.createElement('div');
    darkLayer.id = 'shadow';
    document.body.appendChild(darkLayer);

    const showBlock = document.getElementById('popupWin');
    showBlock.style.display = 'block';

    darkLayer.onclick = () => {
        darkLayer.parentNode.removeChild(darkLayer);
        showBlock.style.display = 'none';
        return false;
    };

    const closeModalContentBtn = document.getElementById("closeModalContentBtn")
    closeModalContentBtn.addEventListener("click", evt =>{
        evt.preventDefault()
        darkLayer.parentNode.removeChild(darkLayer);
        showBlock.style.display = 'none';
    })
}


function setShowModalCreatePost() {
    const addPostBtn = document.getElementById("addPostBtn");
    addPostBtn.addEventListener('click', (evt) => {
        showModalCreatePost();
    })
}



function createLogin() {

    history.pushState(null, null, "/singin");
    document.title = "Login";

    const root = document.getElementById("root")
    root.innerHTML  = `
    <div class="text-center col-2 regAndLogForms">
    <form class="form-signin">
      <img class="mb-4" src="https://проконференции.рф/wp-content/uploads/2017/10/Gerb_MGTU_imeni_Baumana-wpcf_254x300.png" alt="" width="72" height="72">
      <h1 class="h3 mb-3 font-weight-normal">Please sign in</h1>
      <label for="inputEmail" class="sr-only">Email address</label>
      <input type="email" id="inputEmail" class="form-control mb-2" placeholder="Email address">
      <label for="inputPassword" class="sr-only">Password</label>
      <input type="password" id="inputPassword" class="form-control mb-2" placeholder="Password"  pattern=".{6,20}" required title="6 to 20 characters">

      <button id="singinBtn" class="btn btn-lg btn-primary btn-block mb-4" type="submit">Sign in</button>
      <a id="signUpLinkFromLogin"  href="/signup"> <h6 class="linkToAnotherForm"> Sign up</h6> </a>
     
    </form>
     </div>
`
   const signUpLinkFromLogin = document.getElementById("signUpLinkFromLogin")
    signUpLinkFromLogin.addEventListener("click", evt => {
        evt.preventDefault()
        createRegistration()
    })

    const singinBtn = document.getElementById("singinBtn")
    singinBtn.addEventListener("click", evt => {
        evt.preventDefault()

        let inputEmail = document.getElementById("inputEmail").value
        let inputPass = document.getElementById("inputPassword").value


        FetchModule.fetchRequest({url: "/api/signin", body:{email:inputEmail, password:inputPass}})
            .then((res) => {
                return res.ok ? res : Promise.reject(res);
            }).then((response) => {
            return response.json();
        }).then((result) => {



            CurrentUser.Data.login = result.name
            CurrentUser.Data.email = result.email
            CurrentUser.Data.token = result.csrf
            if (CurrentUser.Data.login !== "null") {
                createMainPage()
                return
            }
            alert("ошибка входа")


        }).catch((error) => {
            alert('почта-пароль не совпадает');
        });


    })




}

function createRegistration() {


    history.pushState(null, null, "/singup");
    document.title = "Reg";

    const root = document.getElementById("root")
    root.innerHTML = `
    <div class="text-center col-2 regAndLogForms">
    <form class="form-signin">
      <img class="mb-4" src="https://проконференции.рф/wp-content/uploads/2017/10/Gerb_MGTU_imeni_Baumana-wpcf_254x300.png" alt="" width="72" height="72">
      <h1 class="h3 mb-3 font-weight-normal">Please sign up</h1>
      <input type="text" id="inputName" class="form-control mb-2" placeholder="Nickname" required="" autofocus="">
      <label for="inputEmail" class="sr-only">Email address</label>
      <input type="email" id="inputEmail" class="form-control mb-2" placeholder="Email address" required="" autofocus="">
      <label for="inputPassword" class="sr-only">Password</label>
      <input type="password" id="inputPassword" class="form-control mb-2" placeholder="Password"  pattern=".{6,20}" required title="6 to 20 characters">
      <input type="password" id="inputPassword2" class="form-control mb-2" placeholder="Rewrite password"  pattern=".{6,20}" required title="6 to 20 characters">
      <div class="checkbox mb-3">
      </div>
      <button id="signInBtn" class="btn btn-lg btn-primary btn-block mb-4" type="submit">Sign up</button>
       <a id="signInLinkFromRegistration"  href="/signin"> <h6> Sign in</h6> </a>
     
    </form>
    </div>`

    const signInLinkFromRegistration = document.getElementById("signInLinkFromRegistration")
    signInLinkFromRegistration.addEventListener("click", evt => {
        evt.preventDefault()
        createLogin()
    })

    const signInBtn =document.getElementById("signInBtn")
    signInBtn.addEventListener("click", evt => {

        evt.preventDefault()

        let inputNickname = document.getElementById("inputName").value
        let inputEmail = document.getElementById("inputEmail").value
        let inputPass = document.getElementById("inputPassword").value
        let inputPass2 = document.getElementById("inputPassword2").value

        if ( !(inputPass === inputPass2 && inputPass.length > 6 && inputPass.length < 20) ) {
            alert("wrong rewrite pass and pass length must be [6,20]")
            return
        }

        console.log("send!")

        FetchModule.fetchRequest({url: "/api/signup", body:{name:inputNickname,email:inputEmail, password:inputPass}})
            .then((res) => {
                return res.ok ? res : Promise.reject(res);
            }).then((response) => {
            return response.json();
        }).then((result) => {

            CurrentUser.Data.login = result.name
            CurrentUser.Data.email = result.email
            CurrentUser.Data.token = result.csrf
            createMainPage()

        }).catch((error) => {
            alert('Ошибка (почта-логин занят, некорректные данные)');
        });


    })


}


function createMainPage() {

    history.pushState(null, null, "/main");
    document.title = "Main";

    const root = document.getElementById("root")
    root.innerHTML = ` 
      <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
      <div class="d-flex m-auto">
        <a class="navbar-brand mr-5" style="font-size: 25px;">` + CurrentUser.Data.login + `</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarCollapse">
          <ul class="navbar-nav mr-auto mr-5">
            <li class="nav-item active mr-5">
              <a class="nav-link" style="margin-right: 200px; font-size: 17px"> Email: ` + CurrentUser.Data.email + `<span class="sr-only">(current)</span></a>
           
          </ul>
          <form class="form-inline mt-2 mt-md-0 ml-5">
             <button id="logOutBtn" class="btn btn-outline-success my-2 my-sm-0 ml-5" type="submit">Log out</button>
             
             <button class="btn btn-outline-success my-2 my-sm-0 ml-5"> <a href="/getbadpage">  Bad Page with post titles</a> </button> 
          </form>
        </div>
        </div>
      </nav>
      
    <div id="content" class=" m-auto">
    
    <div class=" " style="text-align: center"> 
    
    <div id="postsSection" class="col-5  ml-auto mr-auto d-inline-block postSection">
  
    </div>
    
    
    <div id="addPostSection" class="d-inline-block " style="box-sizing: content-box; width: 50px; height: 50px; vertical-align: top; margin-top: 150px; ">
    <img id="addPostBtn" class="  mb-4 " src="/plusImg.svg" style="width: 50px; object-fit: fill; cursor: pointer;" >
    </div>
    
   
    </div>
    
    <div id="popupWin" class="modalWindowCreatePost col-5">
    <h3 style="display: inline"> Create post</h3>
    <img src="/closeImg.svg" id="closeModalContentBtn" class="closeImgInModal">   
    
    <div class="input-group mb-3">
  <div class="input-group-prepend">
    <span class="input-group-text" id="inputGroup-sizing-default">Оглавние</span>
  </div>
  <input type="text" id="inputPostTitle" class="form-control" aria-label="Default" aria-describedby="inputGroup-sizing-default" style="height: 20%">
</div>

    <div class="input-group mb-3">
    <div class="input-group-prepend">
        <span class="input-group-text" id="inputGroup-sizing-default">Текст поста</span>
    </div>
    <textarea rows="4"  id="inputPostText" class="form-control" aria-label="With textarea" style="height: 60%; max-height: 60% ; resize: none;"></textarea>
    </div> 
    
    <button id="createPostBtn" class="btn btn-outline-success"> Create </button>
    
    </div>
    
    </div>
    `

    getPosts()
    setShowModalCreatePost()
    setCreatePostBtn();


    const logOutBtn = document.getElementById("logOutBtn")
    logOutBtn.addEventListener("click", evt=>{
        evt.preventDefault()
        FetchModule.fetchRequest({url:"/api/logout"})
            .then((res) => {
                return res.ok ? res : Promise.reject(res);
            }).then((response) => {
            return response.json();
        }).then((result) => {

            createLogin()

        }).catch((error) => {
            alert('Ошибка при выходе:' + error);
        });
    })


}

function setCreatePostBtn() {
    const createPostBtn = document.getElementById("createPostBtn")
    createPostBtn.addEventListener("click", evt=>{
        evt.preventDefault();


        const inputTitleF = document.getElementById("inputPostTitle")
        const inputTextF = document.getElementById("inputPostText")

        const inputTitle = inputTitleF.value
        const inputText = inputTextF.value

        if ( !(inputTitle.length >= 5 && inputTitle.length <= 300 && inputText.length >= 5 && inputText.length < 1000) ) {
            alert("post title length must be [5,300], text length [5,1000]")
            return;
        }

        FetchModule.fetchRequest({url:"/api/createpost", body:{title:inputTitle, text:inputText}})
        .then((res) => {
                return res.ok ? res : Promise.reject(res);
            }).then((response) => {
            return response.json();
        }).then((result) => {

            addPost(result)

            const darkLayer = document.getElementById('shadow');
            const showBlock = document.getElementById('popupWin');
            darkLayer.parentNode.removeChild(darkLayer);
            showBlock.style.display = 'none';
            inputTitleF.value = ""
            inputTextF.value = ""

        }).catch((error) => {
            alert('Ошибка при создании поста. (возможно пост с таким оглавлением уже существует)');
        });


    })

}

function addPost(post) {

    post.date = post.date.replace("T", "   ")
    post.date = post.date.substring(0, 18)

    const postSection = document.getElementById("postsSection")
    let element = document.createElement("div", );

    element.innerHTML = `          <div class="card flex-md-row  mb-4 box-shadow h-md-250 onePost" style="box-shadow: 0 0 5px 2px #878787; text-align: left;">
            <div class="card-body d-flex flex-column align-items-start">
              <h3 class="mb-0">
                <a class="text-dark" >` + post.title + `</a>
              </h3>
              <div style="display: inline-block" > 
              <strong class="d-inline-block mb-2 text-success d-inline mr-2"> ` + post.author + `</strong>
              <div class="mb-1 text-muted d-inline ml-3">`+ post.date + `  </div>
              </div>
              <br> 
              <p class="card-text mb-auto ">` + post.text + `</p>
            </div>
            <img class=" col-2 mt-2 mb-2 card-img-right flex-auto d-none d-md-block mr-4" src="https://проконференции.рф/wp-content/uploads/2017/10/Gerb_MGTU_imeni_Baumana-wpcf_254x300.png" style="width: 80px; object-fit: contain " >
          </div>`
    postSection.appendChild(element)
}


function getPosts() {

    FetchModule.fetchRequest({url:"/api/getposts", method:"get"})
        .then((res) => {
            return res.ok ? res : Promise.reject(res);
        }).then((response) => {
        return response.json();
    }).then((result) => {

        result.forEach((item, i, arr) => {
            addPost(item)
        });

    }).catch((error) => {
       // alert('Ошибка при получении постов:' + error);
    });


}

function start(){
    if (!navigator.onLine) {
        const root = document.getElementById("root")
        root.innerHTML = "<h1> No connection! (global net) </h1>"
        return
    }

    let currentLocation = window.location.pathname;

    FetchModule.fetchRequest({url: "http://127.0.0.1:8080" + "/api/getme"})
        .then((res) => {
            return res.ok ? res : Promise.reject(res);
        }).then((response) => {
        return response.json();
        }).then((result) => {

        CurrentUser.Data.login = result.name
        CurrentUser.Data.email = result.email
        CurrentUser.Data.token = result.csrf
        if (CurrentUser.Data.login === undefined || CurrentUser.Data.login === "null") {
            if (currentLocation === "/singup") {
                createRegistration()
                return;
            }
            createLogin()
            return
        }

        createMainPage();

    }).catch((error) => {
        alert('Ошибка обработки запроса получения текущего пользователя:' + error);
    });
}




start();

