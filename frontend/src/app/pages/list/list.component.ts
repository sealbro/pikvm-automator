import {Component, OnInit} from '@angular/core';
import {ApiModule, Command, PiKvmAutomatorService} from "../../api";
import {MatListModule} from "@angular/material/list";
import {NgForOf} from "@angular/common";

@Component({
  selector: 'app-list',
  standalone: true,
  imports: [ApiModule, MatListModule, NgForOf],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss'
})
export class ListComponent implements OnInit {
  commands: Command[] = [];

  constructor(private piKvmAutomatorService: PiKvmAutomatorService) {
  }

  ngOnInit(): void {
    this.piKvmAutomatorService.piKvmAutomatorCommandList().subscribe(data => {
      if (!data.commands) {
        return;
      }
      this.commands = data.commands;
    });
  }
}
