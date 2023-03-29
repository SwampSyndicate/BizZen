import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FindClassesComponent } from './find-classes.component';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";

describe('FindClassesComponent', () => {
  let component: FindClassesComponent;
  let fixture: ComponentFixture<FindClassesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FindClassesComponent ],
      imports: [
        FormsModule,
        ReactiveFormsModule
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FindClassesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
