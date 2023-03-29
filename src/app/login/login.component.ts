import {Component} from '@angular/core';
import {User} from '../user';
import {UserService} from '../user.service';

//passes info between components
import {Router, ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']

})

export class LoginComponent {


  model = new User("12345","", "", "business", []);

  checkbox : boolean = false;


  constructor(private router:Router, private activatedRoute:ActivatedRoute, private userService:UserService){}

  submitted = false;

  onSubmit() {
    this.submitted = true;


    if(this.checkbox)
    {
      this.model.accountType = "business";
    }
    else
    {
      this.model.accountType = "user";
    }
    
    this.router.navigateByUrl('/profile', {state: {user: this.model }});

  }
  routeToRegister() {
    this.router.navigate(['/register']);
  }
  routeToHome() {
    this.router.navigate(['/']);
  }
}
