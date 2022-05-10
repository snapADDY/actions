import itertools
import json
import os
from typing import IO

import requests
from bandit.core import constants, docs_utils, test_properties

BANDIT_COMMENT_ON_PULL_REQUEST = os.environ.get("BANDIT_COMMENT_ON_PULL_REQUEST")


def comment_on_pull_request(message: str):
    """Comments a message within a pull request."""
    token = os.getenv("INPUT_GITHUB_TOKEN")
    if not token:
        print(message)
        return

    if os.getenv("GITHUB_EVENT_NAME") == "pull_request":
        with open(os.getenv("GITHUB_EVENT_PATH")) as json_file:
            event = json.load(json_file)
            headers_dict = {
                "Accept": "application/vnd.github.v3+json",
                "Authorization": f"token {token}",
            }

            request_path = (
                f"https://api.github.com/repos/{event['repository']['full_name']}/"
                f"issues/{event['number']}/comments"
            )
            requests.post(request_path, headers=headers_dict, json={"body": message})


def generate_metrics_table(manager: "bandit.core.manager.BanditManager") -> str:
    """Generate a metrics table.

    Parameters
    ----------
    manager : bandit.core.manager.BanditManager
        A bandit manager object.

    Returns
    -------
    str
        Part of the bandit report with the severity/confidence matrix as table.
    """
    bits = []

    confidence_mapping = {
        "UNDEFINED": "Undefined confidence",
        "LOW": "Low confidence",
        "MEDIUM": "Medium confidence",
        "HIGH": "High confidence",
    }
    severity_mapping = {
        "UNDEFINED": "Undefined severity",
        "LOW": "Low severity",
        "MEDIUM": "Medium severity",
        "HIGH": "High severity",
    }

    table = get_table(manager)

    # headline
    bits.append("\n### Severity/Confidence matrix\n")

    # table header
    bits.append(f"| | {'|'.join([severity_mapping[c] for c in constants.RANKING])}|")
    bits.append("|-:|-:|-:|-:|-:|")

    # create rows
    for constant in constants.RANKING:
        undefined = table[f"{constant}_UNDEFINED"]
        low = table[f"{constant}_LOW"]
        medium = table[f"{constant}_MEDIUM"]
        high = table[f"{constant}_HIGH"]
        bits.append(
            f"|{confidence_mapping[constant]}|{undefined}|{low}|{medium}|{high}|"
        )

    return "\n".join([bit for bit in bits])


def get_issue_description(
    issue: "bandit.core.issue.Issue",
    indent: str,
    show_line_number: bool = True,
    show_code: bool = True,
    n_lines: int = -1,
) -> str:
    """Returns details of issue within a description. This description contains:
    - severity
    - confidence
    - location
    - more info
    - code snippet

    Parameters
    ----------
    issue : bandit.core.issue.Issue
        Bandit issue object.
    indent : str
        Indentation.
    show_line_number : bool, optional
        Show line number of the issue, by default True
    show_code : bool, optional
        Show code environment of the problematic part, by default True
    n_lines : int, optional
        Number of lines to report. -1 means all, by default -1

    Returns
    -------
    str
        Description of an issue.
    """
    bits = []
    bits.append("<details>")
    bits.append(
        f"<summary><strong>[{issue.test_id}:{issue.test}]</strong> {issue.text}</summary>\n<br>\n"
    )

    bits.append(
        f"|<strong>Severity</strong>| {issue.severity.capitalize()} |\n|:-:|:-:|\n|"
        "<strong>Confidence</strong>| {issue.confidence.capitalize()} |"
    )

    bits.append(
        f"|<strong>Location<strong>| {issue.fname}:{issue.lineno if show_line_number else ''}:{''} |"
    )

    bits.append(f"|<strong>More Info<strong>| {docs_utils.get_url(issue.test_id)} |\n")

    if show_code:
        bits.append("<br>\n\n```python")
        bits.extend(
            [indent + line for line in issue.get_code(n_lines, True).split("\n")]
        )
        bits.append("```\n")

    bits.append("</details>")
    return "\n".join([bit for bit in bits])


def get_table(manager: "bandit.core.manager.BanditManager") -> dict:
    """Returns a 5x5 table with all possible RANKING constants.

    Parameters
    ----------
    manager : bandit.core.manager.BanditManager
        A bandit manager object.

    Returns
    -------
    dict
        Dictionary with all possible combinations of the RANKING constants.
    """
    table = {
        "_".join(combination): 0
        for combination in list(itertools.product(constants.RANKING, constants.RANKING))
    }
    for issue_dict in manager.results:
        issue_confidence = issue_dict.as_dict().get("issue_confidence", 0)
        issue_severity = issue_dict.as_dict().get("issue_severity", 0)
        table[f"{issue_confidence}_{issue_severity}"] += 1

    return table


def get_detailed_results(
    manager: "bandit.core.manager.BanditManager",
    severity_level: str,
    confidence_level: str,
    n_lines: int = -1,
) -> str:
    """Get detailed bandit results for every security issue.

    Parameters
    ----------
    manager : bandit.core.manager.BanditManager
       A bandit manager object.
    severity_level : str
        Filtering severity level ('LOW', 'MEDIUM', 'HIGH').
    confidence_level : str
        Filtering confidence level ('LOW', 'MEDIUM', 'HIGH').
    n_lines : int, optional
        Number of lines to report. -1 means all, by default -1

    Returns
    -------
    str
        Part of the bandit report with the detailed issue descriptions.
    """
    bits = []
    issues = manager.get_issue_list(severity_level, confidence_level)
    baseline = not isinstance(issues, list)
    candidate_indent = " " * 10

    if not len(issues):
        return "\tNo issues identified."

    for issue in issues:
        # if not a baseline or only one candidate we know the issue of
        if not baseline or len(issues[issue]) == 1:
            bits.append(get_issue_description(issue, "", lines=n_lines))

        # otherwise show the finding and the candidates
        else:
            bits.append(
                get_issue_description(
                    issue, "", show_line_number=False, show_code=False
                )
            )

            bits.append("\n-- Candidate Issues --")
            for candidate in issues[issue]:
                bits.append(
                    get_issue_description(candidate, candidate_indent, lines=n_lines)
                )
                bits.append("\n")
    return "\n".join([bit for bit in bits])


def get_verbose_details(manager: "bandit.core.manager.BanditManager"):
    """Get verbose details."""
    bits = []
    bits.append(f"Files in scope ({len(manager.files_list)}):")
    tpl = "\t%s (score: {SEVERITY: %i, CONFIDENCE: %i})"
    bits.extend(
        [
            tpl % (item, sum(score["SEVERITY"]), sum(score["CONFIDENCE"]))
            for (item, score) in zip(manager.files_list, manager.scores)
        ]
    )
    bits.append(f"Files excluded ({len(manager.excluded_files)}):")
    bits.extend(["\t%s" % fname for fname in manager.excluded_files])
    return "\n".join([bit for bit in bits])


@test_properties.accepts_baseline
def report(
    manager: "bandit.core.manager.BanditManager",
    fileobj: IO,
    severity_level: str,
    confidence_level,
    n_lines: int = -1,
):
    """Prints discovered issues in the text format

    Notes
    -----
    Function template taken from:
    https://bandit.readthedocs.io/en/latest/formatters/index.html#example-formatter

    `fileobj` is unused here but required by bandit.

    Parameters
    ----------
    manager : bandit.core.manager.BanditManager
       A bandit manager object.
    fileobj : IO
        A file object.
    severity_level : str
        Filtering severity level ('LOW', 'MEDIUM', 'HIGH').
    confidence_level : str
        Filtering confidence level ('LOW', 'MEDIUM', 'HIGH').
    n_lines : int, optional
        Number of lines to report. -1 means all, by default -1.
    """
    bits = []
    if manager.results_count(severity_level, confidence_level):
        if manager.verbose:
            bits.append(get_verbose_details(manager))

        # add general results
        bits.append("## Bandit Security Check")
        bits.append(
            "<strong>Total lines of code inspected:</strong> "
            f"{(manager.metrics.data['_totals']['loc'])}"
        )
        bits.append(
            "<strong>Total lines skipped (#nosec):</strong> "
            f"{(manager.metrics.data['_totals']['nosec'])}"
        )
        # add metrics table
        bits.append(generate_metrics_table(manager))

        # add detailed security issue results
        bits.append(
            "<details><summary>📋 Click here to see the all possible security issues</summary>\n"
            "<br>\n"
        )
        bits.append(
            get_detailed_results(manager, severity_level, confidence_level, n_lines)
        )
        bits.append("</details>")

        result = "\n".join([bit for bit in bits]) + "\n"

        if BANDIT_COMMENT_ON_PULL_REQUEST:
            comment_on_pull_request(result)
